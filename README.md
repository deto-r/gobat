# gobat
gobat is a library that makes it easy to schedule one-way or parallel batch processes.

## Installing
`$ go get -u github.com/deepoil/gobat`

## Example
**When you want to schedule a one-way batch process...**
```
    import (
        "fmt"
        "gobat"
        "log"
        "time"
    )

    func bat1() {
        fmt.Println("-----start bat1")
        time.Sleep(1000 * time.Millisecond)
        fmt.Println("-----end bat1")
    }
    
    func bat2() {
        fmt.Println("-----start bat2")
        time.Sleep(1000 * time.Millisecond)
        fmt.Println("-----end bat2")
    }

    func main() {
    	// set common config
    	oneWayCommon := gobat.SetCommonBatConfig(
    		time.Now(),
    		1000*10*time.Millisecond,
    		1000*60*60*24*time.Millisecond)

        // one way batch config
	oneWayBat := gobat.SetOneWayBatConfig(oneWayCommon)
        
        // Execution
        for {
            if err := oneWayBat.OneWayBat(bat1, bat2); err != nil {
            	log.Fatal("batch error")
            }
            
            // Set next run time
            oneWayCommon.NextSchedule()
            fmt.Printf("next schedule is: %v\n", oneWayCommon.StartTime)
        }
        
    }
```

**When you want to schedule parallel batch processing...**
```
    import (
        "fmt"
        "gobat"
        "log"
        "time"
    )

    func bat1() {
        fmt.Println("-----start bat1")
        time.Sleep(1000 * time.Millisecond)
        fmt.Println("-----end bat1")
    }
    
    func bat2() {
        fmt.Println("-----start bat2")
        time.Sleep(1000 * time.Millisecond)
        fmt.Println("-----end bat2")
    }

    func bat3() {
        fmt.Println("-----start bat3")
        time.Sleep(1000 * 10 * time.Millisecond)
        fmt.Println("-----end bat3")
    }
    
    func bat4() {
        fmt.Println("-----start bat4")
        time.Sleep(2000 * time.Millisecond)
        fmt.Println("-----end bat4")
    }

    func main() {
    	// set common config
    	paraCommon := gobat.SetCommonBatConfig(
        	time.Now(),
        	1000*10*time.Millisecond,
        	1000*60*60*24*time.Millisecond)

        // Prioritize parallel batches
        p1 := gobat.SetPriority(1, bat1, bat2)
        p2 := gobat.SetPriority(2, bat3)
        p3 := gobat.SetPriority(3, bat4)

        // Set the execution order
	dependency, err := gobat.GenerateDependency(p1, p2, p3)
	if err != nil {
	   log.Fatal(err)
	}

        // Set Config
        paraBat := gobat.SetParallelBatConfig(paraCommon, dependency)

        // Execution
        for {
            if err := paraBat.ParallelBat(); err != nil {
            	log.Fatal("batch error")
            }

            // Set next run time
            paraCommon.NextSchedule()
            fmt.Printf("next schedule is: %v\n", paraCommon.StartTime)
        }
    }
```

# License
MIT License

