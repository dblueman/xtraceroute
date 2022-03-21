package main

import (
   "fmt"
   "os"
   "time"

   "github.com/aeden/traceroute"
   "github.com/dblueman/pinger"
)

func main() {
   if len(os.Args) != 2 {
      fmt.Fprintln(os.Stderr, "usage: xtraceroute <target>")
      os.Exit(1)
   }

   target := os.Args[1]

   p, err := pinger.NewPinger(target, 1000 * time.Millisecond)
   if err != nil {
      fmt.Fprintln(os.Stderr, "NewPinger: "+err.Error())
      os.Exit(1)
   }

   for {
      time.Sleep(60 * time.Second)

      up, err := p.Ping()
      if err != nil {
         fmt.Fprintln(os.Stderr, "Ping: "+err.Error())
         os.Exit(1)
      }

      if up {
         continue
      }

      out, err := traceroute.Traceroute(target, &traceroute.TracerouteOptions{})
      if err != nil {
         fmt.Fprintln(os.Stderr, "Traceroute: "+err.Error())
         os.Exit(1)
      }

      fmt.Println("\nPing timeout at ", time.Now().Format(time.RFC1123Z))

      for _, hop := range out.Hops {
         fmt.Printf("%-3d %v (%v)  %v\n", hop.TTL, hop.HostOrAddressString(), hop.AddressString(), hop.ElapsedTime)
      }
   }
}
