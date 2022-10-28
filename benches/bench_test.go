package benches

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sonyamoonglade/authio/internal/gcmcrypt"
	"github.com/sonyamoonglade/authio/internal/hash"
)

func BenchmarkSign(b *testing.B) {
	key := gcmcrypt.KeyFromString("asdfkjdfkjakafjkjkajkj")
	value := uuid.NewString()
	for i := 0; i < b.N; i++ {
		encrypted, err := gcmcrypt.Encrypt(key, value)
		if err != nil {
			b.Fatalf("err: %v\n", err)
		}
		decrypted, err := gcmcrypt.Decrypt(key, encrypted)
		if value != decrypted {
			b.Fatalf("non eq: %q != %q\n", value, decrypted)
		}
		if err != nil {
			b.Fatalf("err: %v\n", err)
		}
	}

	b.ReportAllocs()
}

func BenchmarkHash(b *testing.B) {

	value := uuid.NewString()
	for i := 0; i < b.N; i++ {
		v := hash.SHA1(value)
		if hash.SHA1(value) != v {
			b.Fail()
		}
	}

	b.ReportAllocs()

}

func TestTimeHashVsSign(t *testing.T) {
	var hashElapsed, signElapsed int64
	value := uuid.NewString()
	wg := new(sync.WaitGroup)
	key := gcmcrypt.KeyFromString("asdfjsdhajahsjdhsjhasdfsadf")
	N := 10_000_0
	fatal := make(chan string, 1)

	wg.Add(1)
	time.Sleep(time.Millisecond)
	go func() {
		start := time.Now()
		for i := 0; i < N; i++ {
			encrypted, err := gcmcrypt.Encrypt(key, value)
			if err != nil {
				fatal <- fmt.Sprintf("err: %v\n", err)
				return
			}
			decrypted, err := gcmcrypt.Decrypt(key, encrypted)
			if err != nil {
				fatal <- fmt.Sprintf("err: %v\n", err)
				return
			}

			if decrypted != value {
				fatal <- fmt.Sprintf("non equal: %q != %q\n", value, decrypted)
				return
			}

		}
		end := time.Now().Sub(start).Milliseconds()
		signElapsed = end
		t.Logf("signing elapsed: %dms\n", end)
		fatal <- "OK"
		defer wg.Done()
	}()

	wg.Add(1)
	time.Sleep(time.Millisecond)
	go func() {
		start := time.Now()
		for i := 0; i < N; i++ {
			h := hash.SHA1(value)
			if hash.SHA1(value) != h {
				fatal <- fmt.Sprintf("not equal: %q != %q\n", h, hash.SHA1(value))
				return
			}
		}
		end := time.Now().Sub(start).Milliseconds()
		t.Logf("hashing elapsed: %dms\n", end)
		hashElapsed = end
		fatal <- "OK"
		defer wg.Done()
	}()

	f := <-fatal
	if f != "OK" {
		t.Fatal(f)
	}

	wg.Wait()

	t.Logf("signing is %.2fx times slower than hashing with N=%d", float64(signElapsed)/float64(hashElapsed), N)
}

//todo: rps test
