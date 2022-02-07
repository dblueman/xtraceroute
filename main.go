package main

import (
   "fmt"
   "os"
   "time"

   "github.com/aeden/traceroute"
)

func dump(out traceroute.TracerouteResult) {
   t := time.Now().Format(time.RFC1123Z)
   fmt.Printf(" at %s\n", t)

   for _, hop := range out.Hops {
      fmt.Printf("%-3d %v (%v)  %v\n", hop.TTL, hop.HostOrAddressString(), hop.AddressString(), hop.ElapsedTime)
   }

   fmt.Println()
}

func diff(current, last traceroute.TracerouteResult) bool {
   if len(current.Hops) != len(last.Hops) {
      fmt.Printf("number of hops went from %d to %d", len(last.Hops), len(current.Hops))
      return true
   }

   for i := range current.Hops {
      if current.Hops[i].AddressString() != last.Hops[i].AddressString() {
         fmt.Printf("hop %d changed from %s to %s", i, last.Hops[i].AddressString(), current.Hops[i].AddressString())
         return true
      }
   }

   return false
}

func main() {
   if len(os.Args) != 2 {
      fmt.Fprintln(os.Stderr, "usage: xtraceroute <target>")
      os.Exit(1)
   }

   target := os.Args[1]
   var last traceroute.TracerouteResult

   for {
      out, err := traceroute.Traceroute(target, &traceroute.TracerouteOptions{})
      if err != nil {
         fmt.Fprintln(os.Stderr, "failed: "+err.Error())
         os.Exit(1)
      }

      if diff(out, last) {
         dump(out)
         last = out
      }

      time.Sleep(60 * time.Second)
   }
}
