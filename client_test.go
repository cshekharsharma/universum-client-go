package universum

import (
	"context"
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestNewClient_Success(t *testing.T) {
	opts := &Options{
		ConnPoolsize: 10,
	}

	client, err := NewClient(opts)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.pool == nil {
		t.Fatal("Expected connection pool to be initialized, got nil")
	}

	if client.opts != opts {
		t.Fatalf("Expected client opts to match input, got %v", client.opts)
	}

	if client.id == "" {
		t.Fatal("Expected client ID to be set, but got an empty ID")
	}

	decodedId, err := base64.StdEncoding.DecodeString(client.id)
	if err != nil {
		t.Fatalf("Expected valid base64 ID, but decoding failed: %v", err)
	}

	_, err = strconv.Atoi(string(decodedId))
	if err != nil {
		t.Fatalf("Expected client ID to be a valid timestamp, but got: %v", err)
	}
}

func TestClient_Commands(t *testing.T) {
	opts := mockOptions()

	client, err := NewClient(opts)
	if err != nil {
		t.Fatalf("Expected no error while creating client, got %v", err)
	}

	ctx := context.Background()
	intKey := fmt.Sprintf("int-%v", time.Now().UnixNano())
	stringKey := fmt.Sprintf("string-%v", time.Now().UnixNano())
	listKey := fmt.Sprintf("list-%v", time.Now().UnixNano())

	pingResult, err := client.Ping(ctx)
	if err != nil {
		t.Fatalf("Ping:: Expected no error from Get, got %v", err)
	} else if err = isSuccessResponse(pingResult.Code, pingResult.Message, RespPingSuccess, "OK"); err != nil {
		t.Fatal("Ping: " + err.Error())
	}

	// Get non-existant key
	getResult, err := client.Get(ctx, intKey)
	if err != nil {
		t.Fatalf("GET:: Expected no error from Get, got %v", err)
	} else if err = isSuccessResponse(getResult.Code, getResult.Value, RespRecordNotFound, nil); err != nil {
		t.Fatal("Get: " + err.Error())
	}

	// Exists for non-existant key
	existsResult, err := client.Exists(ctx, stringKey)
	if err != nil {
		t.Fatalf("EXISTS:: Expected no error from Get, got %v", err)
	} else if err = isSuccessResponse(existsResult.Code, existsResult.Found, RespRecordNotFound, false); err != nil {
		t.Fatal("Exists: " + err.Error())
	}

	// Set a key
	setResult, err := client.Set(ctx, stringKey, "prefix_", 5)
	if err != nil {
		t.Fatalf("Set:: Expected no error from Get, got %v", err)
	} else if err = isSuccessResponse(setResult.Code, setResult.Success, RespRecordUpdated, true); err != nil {
		t.Fatal("Set: " + err.Error())
	}

	time.Sleep(1 * time.Second) // Wait for a second to check TTL correctness

	// TTL for existing key
	ttlResult, err := client.TTL(ctx, stringKey)
	if err != nil {
		t.Fatalf("TTL:: Expected no error from Get, got %v", err)
	} else if err = isSuccessResponse(ttlResult.Code, ttlResult.TTL, RespRecordFound, 4*time.Second); err != nil {
		t.Fatal("TTL: " + err.Error())
	}

	// Expire existing key
	expireResult, err := client.Expire(ctx, stringKey, 8)
	if err != nil {
		t.Fatalf("Expire:: Expected no error from Get, got %v", err)
	} else if err = isSuccessResponse(expireResult.Code, expireResult.Success, RespRecordUpdated, true); err != nil {
		t.Fatal("Expire: " + err.Error())
	}

	// Append existing key
	append2Result, err := client.Append(ctx, stringKey, "suffix")
	if err != nil {
		t.Fatalf("Append2:: Expected no error from Get, got %v", err)
	} else if err = isSuccessResponse(append2Result.Code, append2Result.ContentLength, RespRecordUpdated, int64(13)); err != nil {
		t.Fatal("Append2: " + err.Error())
	}

	// Get existant key
	getResult2, err := client.Get(ctx, stringKey)
	if err != nil {
		t.Fatalf("Get2:: Expected no error from Get, got %v", err)
	} else if err = isSuccessResponse(getResult2.Code, getResult2.Value, RespRecordFound, "prefix_suffix"); err != nil {
		t.Fatal("Get2: " + err.Error())
	}

	time.Sleep(6 * time.Second) // Wait for a few seconds

	// Exists for existant key
	exists2Result, err := client.Exists(ctx, stringKey)
	if err != nil {
		t.Fatalf("Exists2:: Expected no error from Get, got %v", err)
	} else if err = isSuccessResponse(exists2Result.Code, exists2Result.Found, RespRecordFound, true); err != nil {
		t.Fatal("Exists2: " + err.Error())
	}

	time.Sleep(4 * time.Second) // Wait for key to expire

	// Append non-existing key
	append1Result, err := client.Append(ctx, stringKey, "NextSuffix")
	if err != nil {
		t.Fatalf("Append1:: Expected no error from Get, got %v", err)
	} else if err = isSuccessResponse(append1Result.Code, append1Result.ContentLength, RespRecordNotFound, int64(-99999999)); err != nil {
		t.Fatal("Append1: " + err.Error())
	}

	// MSet multiple keys
	msetResult, err := client.MSet(ctx, map[string]interface{}{intKey: 990, listKey: []interface{}{true, 0}})
	if err != nil {
		t.Fatalf("MSet:: Expected no error from Get, got %v", err)
	} else if err = isSuccessResponse(msetResult.Code, msetResult.Successes, RespMsetCompleted, map[string]bool{intKey: true, listKey: true}); err != nil {
		t.Fatal("MSet: " + err.Error())
	}

	// Increment existing key
	incResult, err := client.Increment(ctx, intKey, 9)
	if err != nil {
		t.Fatalf("Increment:: Expected no error from Get, got %v", err)
	} else if err = isSuccessResponse(incResult.Code, incResult.NewValue, RespRecordUpdated, int64(999)); err != nil {
		t.Fatal("Increment: " + err.Error())
	}

	// Decrement existing key
	decrResult, err := client.Decrement(ctx, intKey, 20)
	if err != nil {
		t.Fatalf("Decrement:: Expected no error from Get, got %v", err)
	} else if err = isSuccessResponse(decrResult.Code, decrResult.NewValue, RespRecordUpdated, int64(979)); err != nil {
		t.Fatal("Decrement: " + err.Error())
	}

	// Delete a key
	deleteResult, err := client.Delete(ctx, intKey)
	if err != nil {
		t.Fatalf("Delete:: Expected no error from Get, got %v", err)
	} else if err = isSuccessResponse(deleteResult.Code, deleteResult.Deleted, RespRecordDeleted, true); err != nil {
		t.Fatal("Delete: " + err.Error())
	}

	//MGet multiple keys
	mgetResult, err := client.MGet(ctx, []string{intKey, stringKey, listKey})
	if err != nil {
		t.Fatalf("MGet:: Expected no error from Get, got %v", err)
	} else if err = isSuccessResponse(mgetResult.Code, mgetResult.Values, RespMgetCompleted,
		map[string]interface{}{
			intKey:    map[string]interface{}{"Code": RespRecordNotFound, "Value": interface{}(nil)},
			stringKey: map[string]interface{}{"Code": RespRecordNotFound, "Value": interface{}(nil)},
			listKey:   map[string]interface{}{"Code": RespRecordFound, "Value": []interface{}{true, int64(0)}},
		}); err != nil {
		t.Fatal("MGet: " + err.Error())
	}

	// MDelete multiple keys
	mdelResult, err := client.MDelete(ctx, []string{intKey, stringKey, listKey})
	if err != nil {
		t.Fatalf("MDelete:: Expected no error from Get, got %v", err)
	} else if err = isSuccessResponse(mdelResult.Code, mdelResult.Deletions, RespMdelCompleted,
		map[string]bool{intKey: true, stringKey: true, listKey: true}); err != nil {
		t.Fatal("MDelete: " + err.Error())
	}

	// MGet multiple keys
	mget2Result, err := client.MGet(ctx, []string{intKey, stringKey, listKey})
	if err != nil {
		t.Fatalf("MGet2:: Expected no error from Get, got %v", err)
	} else if err = isSuccessResponse(mget2Result.Code, mget2Result.Values, RespMgetCompleted,
		map[string]interface{}{
			intKey:    map[string]interface{}{"Code": RespRecordNotFound, "Value": interface{}(nil)},
			stringKey: map[string]interface{}{"Code": RespRecordNotFound, "Value": interface{}(nil)},
			listKey:   map[string]interface{}{"Code": RespRecordNotFound, "Value": interface{}(nil)},
		}); err != nil {
		t.Fatal("MGet2: " + err.Error())
	}
}

func isSuccessResponse(code int64, val interface{}, expectedCode int64, expectedVal interface{}) error {
	if code != expectedCode {
		return fmt.Errorf("Expected code to be %d, got %d", expectedCode, code)
	}

	if !reflect.DeepEqual(val, expectedVal) {
		return fmt.Errorf("Expected value to be %#v, got %#v", expectedVal, val)
	}

	return nil
}

// run this test only with a running universum server
func BenchmarkClient_Benchmark(b *testing.B) {
	opts := mockOptions()

	client, _ := NewClient(opts)
	ctx := context.Background()
	var ttl int64 = 180

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("K-%v", time.Now().UnixNano())
		_, _ = client.Exists(ctx, key)
		_, _ = client.Set(ctx, key, 1034, ttl)
		_, _ = client.Get(ctx, key)
		_, _ = client.Increment(ctx, key, 12)
		_, _ = client.Decrement(ctx, key, 22)
		_, _ = client.MSet(ctx, map[string]interface{}{key + "AA": 200, key + "BB": "USD"})
		_, _ = client.Expire(ctx, key+"AA", ttl)
		_, _ = client.Expire(ctx, key+"BB", ttl)
		_, _ = client.MGet(ctx, []string{key + "AA", key + "BB", "CC"})
		//	_, _ = client.MDelete(ctx, []string{key + "AA", key + "BB"})
		_, _ = client.TTL(ctx, key)
		//	_, _ = client.Delete(ctx, key)
		_, _ = client.Ping(ctx)
		//_, _ = client.Info(ctx)
	}
}
