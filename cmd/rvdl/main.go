package main

import "log"

func main() {
	pub := NewPublicClient("rvdl test")
	pri, err := NewPrivateClient("tnoniKJQx53kLQ", "HpRn0Kwl02R8XSsH_sXGk4cEfsM", "redditvideodownload", "7tFTuonnBnu8vD", "rvdl test")
	if err != nil {
		log.Fatalln(err)
	}

	pub.Get()

}
