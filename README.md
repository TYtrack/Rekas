<!--
 * @Author: your name
 * @Date: 2021-12-23 00:55:53
 * @LastEditTime: 2021-12-26 19:54:16
 * @LastEditors: TYtrack
 * @Description: 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 * @FilePath: /Rekas/README.md
-->
# Rekas
### **项目介绍**

> Rekas：一个轻量级分布式缓存系统<del>框架</del> ，解决缓存系统中出现的缓存击穿[锁机制]、缓存穿透[布隆过滤器]、缓存雪崩[分布式]问题，实现了 Lightweight "Remote Cache Access" (Rekas) Framework

### **实现目标**

<br>

- </del>单机缓存和基于 HTTP 的分布式缓存</del>
- <del>最近最少访问(Least Recently Used, LRU) 缓存策略</del>
- <del>实现Master服务器对分布式服务器的管理</del>
- <del>增加心跳检测来实现对分布服务器存活的检测</del>
- <del>使用一致性哈希选择节点，实现负载均衡</del>
- <del>使用Viper库实现配置解析管理</del>
- <del>利用锁机制防止缓存击穿</del>
- 使用布隆过滤器实现缓存穿透【集成中】