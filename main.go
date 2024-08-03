package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	defer func() {
		fmt.Println(time.Since(start))
	}()

	evilNinjas := []string{"Tommy", "Johnny", "Bobby", "Andy"}

	smokeSignal := make(chan bool, 4)
	for _, evilNinja := range evilNinjas {
		go attack(evilNinja, smokeSignal)
	}

	for i, _ := range evilNinjas {
		fmt.Println(<-smokeSignal, i)
	}
	// time.Sleep(time.Second * 2)

	channel := make(chan string)
	go throwingNinjaStar(channel)
	for {
		message, open := <-channel
		if !open {
			break
		}
		fmt.Println(message)
	}

	ninja1, ninja2 := make(chan string), make(chan string)

	go captainElect(ninja1, "Ninja 1")
	go captainElect(ninja2, "Ninja 2")

	// fmt.Println(<-ninja1)
	// fmt.Println(<-ninja2)

	select {
	case message := <-ninja1:
		fmt.Println(message)
	case message := <-ninja2:
		fmt.Println(message)
	default:
		fmt.Println("Neither")

	}

	roughlyFair()

	// time.Sleep(time.Millisecond * 1000)

	var beeper sync.WaitGroup
	evilNinjasBeep := []string{"Tommy", "Johnny", "Bobby"}
	beeper.Add(len(evilNinjasBeep))
	for _, evilNinja := range evilNinjasBeep {
		go attackBeep(evilNinja, &beeper)
	}
	beeper.Wait()
	fmt.Println("Mission completed")

	var beeper2 sync.WaitGroup
	beeper2.Add(1)
	beeper2.Done() // leaving this out or calling more than added is a deadlock
	beeper2.Wait()

	mutexTutorial()

	onceTutorial()

	signalBroadcastTutorial()

}

func signalBroadcastTutorial() {
	gettingReadyForMissionWithCond()
	broadcastStartOfMission()
}

var ready bool

func broadcastStartOfMission() {
	beeper := sync.NewCond(&sync.Mutex{})
	var wg sync.WaitGroup
	wg.Add(3)

	standByForMission(func() {
		fmt.Println("Ninja 1 starting mission.")
		wg.Done()
	}, beeper)

	standByForMission(func() {
		fmt.Println("Ninja 2 starting mission.")
		wg.Done()
	}, beeper)

	standByForMission(func() {
		fmt.Println("Ninja 3 starting mission.")
		wg.Done()
	}, beeper)

	beeper.Broadcast()
	wg.Wait()
	fmt.Println("All Ninjas have started their missions")
}

func standByForMission(fn func(), beeper *sync.Cond) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		beeper.L.Lock()
		defer beeper.L.Unlock()
		beeper.Wait()
		fn()
	}()
	wg.Wait()
}

func gettingReadyForMissionWithCond() {
	cond := sync.NewCond(&sync.Mutex{})

	go gettingReadyWithCond(cond)
	workIntervals := 0
	cond.L.Lock()
	for !ready {
		workIntervals++
		cond.Wait()
	}
	cond.L.Unlock()
	fmt.Printf("We are now ready! After %d work intervals.\n", workIntervals)
}
func gettingReadyForMission() {
	go gettingReady()
	workIntervals := 0
	for !ready {
		time.Sleep(3 * time.Second)
		workIntervals++
	}
	fmt.Printf("We are now ready! After %d work intervals.\n", workIntervals)
}

func gettingReadyWithCond(cond *sync.Cond) {
	sleep()
	ready = true
	cond.Signal()
}

func gettingReady() {
	sleep()
	ready = true
}

func sleep() {
	rand.Seed(time.Now().UnixNano())
	someTime := time.Duration(1+rand.Intn(3)) * time.Second
	time.Sleep(someTime)
}

var missionCompleted bool

func onceTutorial() {

	var wg sync.WaitGroup
	wg.Add(100)

	var once sync.Once

	for i := 0; i < 100; i++ {

		go func() {
			if foundTreasure() {
				// this only calls this method once but the for loop continues
				once.Do(markMissionCompleted)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	checkMissionCompletion()
}

func checkMissionCompletion() {
	if missionCompleted {
		fmt.Println("Misson is now completed.")
	} else {
		fmt.Println("Mission was a failure.")
	}
}

func markMissionCompleted() {
	missionCompleted = true
	fmt.Println("Marking mission completed")
}

func foundTreasure() bool {
	rand.Seed(time.Now().UnixNano())
	return 0 == rand.Intn(10)
}

var (
	lock   sync.Mutex
	rwLock sync.RWMutex
	count  int
)

func mutexTutorial() {
	// mutexBasics()

	readAndWrite()
}

func readAndWrite() {
	go read()
	go read()
	go write()

	time.Sleep(1 * time.Second)
	fmt.Println("Done")
}

func read() {
	rwLock.RLock()
	defer rwLock.RUnlock()
	fmt.Println("Read locking")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Reading unlocking")
}

func write() {
	rwLock.Lock()
	defer rwLock.Unlock()
	fmt.Println("Write locking")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Write unlocking")
}

func mutexBasics() {
	iterations := 100
	for i := 0; i < iterations; i++ {
		go increment()
	}
	time.Sleep(1 * time.Second)
	fmt.Println("Resulted count is:", count)
}

func increment() {
	lock.Lock()
	count++
	lock.Unlock()
}

func attackBeep(evilNinja string, beeper *sync.WaitGroup) {
	fmt.Println("Attacked evil ninja: ", evilNinja)
	beeper.Done()
}

func roughlyFair() {
	ninja1 := make(chan interface{})
	close(ninja1)
	ninja2 := make(chan interface{})
	close(ninja2)

	var ninja1Count, ninja2Count int
	for i := 0; i < 1000; i++ {
		select {
		case <-ninja1:
			ninja1Count++
		case <-ninja2:
			ninja2Count++
		}
	}

	fmt.Printf("ninja1Count: %d, ninja2Count: %d\n", ninja1Count, ninja2Count)
}
func captainElect(ninja chan string, message string) {
	ninja <- message
}

func throwingNinjaStar(channel chan string) {
	rand.Seed(time.Now().UnixNano())
	numRounds := 3
	for i := 0; i < numRounds; i++ {
		score := rand.Intn(10)
		channel <- fmt.Sprint("You scored: ", score)
	}
	close(channel)
}

func attack(target string, attacked chan bool) {
	time.Sleep(time.Millisecond * 10)
	fmt.Println("Throwing ninja stars at", target)
	attacked <- true
}
