package main

import (
	"fmt"
	gobat2 "gobat"
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

func bat5() {
	fmt.Println("-----start bat5")
	time.Sleep(2000 * time.Millisecond)
	fmt.Println("-----end bat5")
}

func bat6() {
	fmt.Println("-----start bat6")
	time.Sleep(2000 * time.Millisecond)
	fmt.Println("-----end bat6")
}

func bat7() {
	fmt.Println("-----start bat7")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("-----end bat7")
}

func main() {
	// set common config
	oneWayCommon := gobat2.SetCommonBatConfig(
		time.Now(),
		1000*10*time.Millisecond,
		1000*60*60*24*time.Millisecond)

	paraCommon := gobat2.SetCommonBatConfig(
		time.Now(),
		1000*10*time.Millisecond,
		1000*60*60*24*time.Millisecond)

	// Prioritize parallel batches
	p1 := gobat2.SetPriority(1, bat3, bat4, bat1, bat2)
	p2 := gobat2.SetPriority(2, bat5)
	p3 := gobat2.SetPriority(3, bat6, bat7)

	// Set the execution order
	dependency, err := gobat2.GenerateDependency(p1, p2, p3)
	if err != nil {
		log.Fatal(err)
	}

	// Set Config
	paraBat := gobat2.SetParallelBatConfig(paraCommon, dependency)
	// one way batch
	oneWayBat := gobat2.SetOneWayBatConfig(oneWayCommon)

	// Execution example
	for {
		if err := paraBat.ParallelBat(); err != nil {
			log.Fatal("batch error")
		}
		paraCommon.NextSchedule()
		fmt.Printf("next schedule is: %v\n", paraCommon.StartTime)

		if err := oneWayBat.OneWayBat(bat1, bat2, bat3, bat4); err != nil {
			log.Fatal("batch error")
		}
		oneWayCommon.NextSchedule()
		fmt.Printf("next schedule is: %v\n", oneWayCommon.StartTime)
	}

}
