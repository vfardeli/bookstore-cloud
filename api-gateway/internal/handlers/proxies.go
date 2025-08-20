package handlers

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"api-gateway/internal/utils"

	"github.com/cenkalti/backoff/v4"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/sony/gobreaker"
	"go.opentelemetry.io/otel/propagation"
)

var (
	bookBreaker  *gobreaker.CircuitBreaker
	orderBreaker *gobreaker.CircuitBreaker
	userBreaker  *gobreaker.CircuitBreaker
	log          = logrus.New()
)

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)

	bookBreaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "BookService",
		MaxRequests: 5,
		Interval:    60 * time.Second, // rolling window for counts
		Timeout:     10 * time.Second, // how long to stay open
	})

	orderBreaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "OrderService",
		MaxRequests: 5,
		Interval:    60 * time.Second, // rolling window for counts
		Timeout:     10 * time.Second, // how long to stay open
	})

	userBreaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "UserService",
		MaxRequests: 5,
		Interval:    60 * time.Second, // rolling window for counts
		Timeout:     10 * time.Second, // how long to stay open
	})
}

func ProxyToBookService(c *gin.Context) {
	reqID := uuid.New().String()

	_, err := bookBreaker.Execute(func() (interface{}, error) {
		client := &http.Client{Timeout: 5 * time.Second}

		// 🔹 Capture the remaining path after /books
		restOfPath := c.Param("path") // includes leading slash if any
		targetURL := "http://book-service:8002/books" + restOfPath

		// Add query params if exist
		if c.Request.URL.RawQuery != "" {
			targetURL += "?" + c.Request.URL.RawQuery
		}

		utils.SendLog("api-gateway", reqID, "info", "Forwarding request", map[string]interface{}{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"target": targetURL,
		})

		// Copy body (for POST/PUT)
		var reqBody []byte
		if c.Request.Body != nil {
			data, _ := io.ReadAll(c.Request.Body)
			defer c.Request.Body.Close()
			reqBody = data
		}

		operation := func() error {
			var bodyReader io.Reader
			if reqBody != nil {
				bodyReader = bytes.NewBuffer(reqBody) // re-use saved body
			}

			req, err := http.NewRequest(c.Request.Method, targetURL, bodyReader)
			if err != nil {
				return backoff.Permanent(err) // don’t retry if request invalid
			}

			// Copy headers
			for k, v := range c.Request.Header {
				req.Header[k] = v
			}
			req.Header.Set("X-Request-ID", reqID)

			// Ensure book-service can continue tracing
			propagator := propagation.TraceContext{}
			propagator.Inject(c.Request.Context(), propagation.HeaderCarrier(req.Header))

			resp, err := client.Do(req)
			if err != nil {
				utils.SendLog("api-gateway", reqID, "warn", "Request failed, retrying", map[string]interface{}{
					"error": err.Error(),
				})
				return err
			}
			defer resp.Body.Close()

			// Only retry on 5xx errors
			if resp.StatusCode >= 500 {
				utils.SendLog("api-gateway", reqID, "warn", "Book service error, retrying", map[string]interface{}{
					"status": resp.StatusCode,
				})
				return io.EOF
			}

			// Write response
			for k, vv := range resp.Header {
				for _, v := range vv {
					c.Writer.Header().Add(k, v)
				}
			}
			c.Writer.Header().Set("X-Request-ID", reqID)
			c.Writer.WriteHeader(resp.StatusCode)
			_, copyErr := io.Copy(c.Writer, resp.Body)
			return copyErr
		}

		// 🔹 Retry with exponential backoff (max 3 retries, capped delay 5s)
		expBackoff := backoff.NewExponentialBackOff()
		expBackoff.MaxElapsedTime = 5 * time.Second

		err := backoff.Retry(operation, expBackoff)
		return nil, err
	})

	if err != nil {
		utils.SendLog("api-gateway", reqID, "error", "Circuit breaker triggered or retries failed", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":      "Book Service unavailable, try later",
			"request_id": reqID,
		})
		return
	}
}

func ProxyToOrderService(c *gin.Context) {
	reqID := uuid.New().String()

	_, err := orderBreaker.Execute(func() (interface{}, error) {
		client := &http.Client{Timeout: 5 * time.Second}

		// 🔹 Capture the remaining path after /orders
		restOfPath := c.Param("path") // includes leading slash if any
		targetURL := "http://order-service:8003/orders" + restOfPath

		// Add query params if exist
		if c.Request.URL.RawQuery != "" {
			targetURL += "?" + c.Request.URL.RawQuery
		}

		utils.SendLog("api-gateway", reqID, "info", "Forwarding request", map[string]interface{}{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"target": targetURL,
		})

		// Copy body (for POST/PUT)
		var reqBody []byte
		if c.Request.Body != nil {
			data, _ := io.ReadAll(c.Request.Body)
			defer c.Request.Body.Close()
			reqBody = data
		}

		operation := func() error {
			var bodyReader io.Reader
			if reqBody != nil {
				bodyReader = bytes.NewBuffer(reqBody) // re-use saved body
			}

			req, err := http.NewRequest(c.Request.Method, targetURL, bodyReader)
			if err != nil {
				return backoff.Permanent(err) // don’t retry if request invalid
			}

			// Copy headers
			for k, v := range c.Request.Header {
				req.Header[k] = v
			}
			req.Header.Set("X-Request-ID", reqID)

			resp, err := client.Do(req)
			if err != nil {
				utils.SendLog("api-gateway", reqID, "warn", "Request failed, retrying", map[string]interface{}{
					"error": err.Error(),
				})
				return err
			}
			defer resp.Body.Close()

			// Only retry on 5xx errors
			if resp.StatusCode >= 500 {
				utils.SendLog("api-gateway", reqID, "warn", "Order service error, retrying", map[string]interface{}{
					"status": resp.StatusCode,
				})
				return io.EOF
			}

			// Write response
			for k, vv := range resp.Header {
				for _, v := range vv {
					c.Writer.Header().Add(k, v)
				}
			}
			c.Writer.Header().Set("X-Request-ID", reqID)
			c.Writer.WriteHeader(resp.StatusCode)
			_, copyErr := io.Copy(c.Writer, resp.Body)
			return copyErr
		}

		// 🔹 Retry with exponential backoff (max 3 retries, capped delay 5s)
		expBackoff := backoff.NewExponentialBackOff()
		expBackoff.MaxElapsedTime = 5 * time.Second

		err := backoff.Retry(operation, expBackoff)
		return nil, err
	})

	if err != nil {
		utils.SendLog("api-gateway", reqID, "error", "Circuit breaker triggered or retries failed", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":      "Order Service unavailable, try later",
			"request_id": reqID,
		})
		return
	}
}

func ProxyToRegisterUserService(c *gin.Context) {
	reqID := uuid.New().String()

	_, err := userBreaker.Execute(func() (interface{}, error) {
		client := &http.Client{Timeout: 5 * time.Second}

		targetURL := "http://user-service:8001/register"

		// Add query params if exist
		if c.Request.URL.RawQuery != "" {
			targetURL += "?" + c.Request.URL.RawQuery
		}

		utils.SendLog("api-gateway", reqID, "info", "Forwarding request", map[string]interface{}{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"target": targetURL,
		})

		// Copy body (for POST/PUT)
		var reqBody []byte
		if c.Request.Body != nil {
			data, _ := io.ReadAll(c.Request.Body)
			defer c.Request.Body.Close()
			reqBody = data
		}

		operation := func() error {
			var bodyReader io.Reader
			if reqBody != nil {
				bodyReader = bytes.NewBuffer(reqBody) // re-use saved body
			}

			req, err := http.NewRequest(c.Request.Method, targetURL, bodyReader)
			if err != nil {
				return backoff.Permanent(err) // don’t retry if request invalid
			}

			// Copy headers
			for k, v := range c.Request.Header {
				req.Header[k] = v
			}
			req.Header.Set("X-Request-ID", reqID)

			resp, err := client.Do(req)
			if err != nil {
				utils.SendLog("api-gateway", reqID, "warn", "Request failed, retrying", map[string]interface{}{
					"error": err.Error(),
				})
				return err
			}
			defer resp.Body.Close()

			// Only retry on 5xx errors
			if resp.StatusCode >= 500 {
				utils.SendLog("api-gateway", reqID, "warn", "User service error, retrying", map[string]interface{}{
					"status": resp.StatusCode,
				})
				return io.EOF
			}

			// Write response
			for k, vv := range resp.Header {
				for _, v := range vv {
					c.Writer.Header().Add(k, v)
				}
			}
			c.Writer.Header().Set("X-Request-ID", reqID)
			c.Writer.WriteHeader(resp.StatusCode)
			_, copyErr := io.Copy(c.Writer, resp.Body)
			return copyErr
		}

		// 🔹 Retry with exponential backoff (max 3 retries, capped delay 5s)
		expBackoff := backoff.NewExponentialBackOff()
		expBackoff.MaxElapsedTime = 5 * time.Second

		err := backoff.Retry(operation, expBackoff)
		return nil, err
	})

	if err != nil {
		utils.SendLog("api-gateway", reqID, "error", "Circuit breaker triggered or retries failed", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":      "User Service unavailable, try later",
			"request_id": reqID,
		})
		return
	}
}

func ProxyToLoginService(c *gin.Context) {
	reqID := uuid.New().String()

	_, err := userBreaker.Execute(func() (interface{}, error) {
		client := &http.Client{Timeout: 5 * time.Second}

		targetURL := "http://user-service:8001/login"

		// Add query params if exist
		if c.Request.URL.RawQuery != "" {
			targetURL += "?" + c.Request.URL.RawQuery
		}

		utils.SendLog("api-gateway", reqID, "info", "Forwarding request", map[string]interface{}{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"target": targetURL,
		})

		// Copy body (for POST/PUT)
		var reqBody []byte
		if c.Request.Body != nil {
			data, _ := io.ReadAll(c.Request.Body)
			defer c.Request.Body.Close()
			reqBody = data
		}

		operation := func() error {
			var bodyReader io.Reader
			if reqBody != nil {
				bodyReader = bytes.NewBuffer(reqBody) // re-use saved body
			}

			req, err := http.NewRequest(c.Request.Method, targetURL, bodyReader)
			if err != nil {
				return backoff.Permanent(err) // don’t retry if request invalid
			}

			// Copy headers
			for k, v := range c.Request.Header {
				req.Header[k] = v
			}
			req.Header.Set("X-Request-ID", reqID)

			resp, err := client.Do(req)
			if err != nil {
				utils.SendLog("api-gateway", reqID, "warn", "Request failed, retrying", map[string]interface{}{
					"error": err.Error(),
				})
				return err
			}
			defer resp.Body.Close()

			// Only retry on 5xx errors
			if resp.StatusCode >= 500 {
				utils.SendLog("api-gateway", reqID, "warn", "User service error, retrying", map[string]interface{}{
					"status": resp.StatusCode,
				})
				return io.EOF
			}

			// Write response
			for k, vv := range resp.Header {
				for _, v := range vv {
					c.Writer.Header().Add(k, v)
				}
			}
			c.Writer.Header().Set("X-Request-ID", reqID)
			c.Writer.WriteHeader(resp.StatusCode)
			_, copyErr := io.Copy(c.Writer, resp.Body)
			return copyErr
		}

		// 🔹 Retry with exponential backoff (max 3 retries, capped delay 5s)
		expBackoff := backoff.NewExponentialBackOff()
		expBackoff.MaxElapsedTime = 5 * time.Second

		err := backoff.Retry(operation, expBackoff)
		return nil, err
	})

	if err != nil {
		utils.SendLog("api-gateway", reqID, "error", "Circuit breaker triggered or retries failed", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":      "User Service unavailable, try later",
			"request_id": reqID,
		})
		return
	}
}
