### view
1. 基于gopacket开发
2. gopacket在producer/random_packet_generator.go中使用并发效果很差，推测原因可能是gopacket发包对并发不友好
3. 目前实现的单线程生成流producer/random_packet_generator_simple.go一次只能发一条，pps在30w左右。