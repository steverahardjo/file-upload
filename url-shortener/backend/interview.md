# Interview Q

## When microservice is too micro?
- When it's too tightly coupled with each other: The front facing features require it multiple to be run in-chain.
- When there is no specific feature or concern that it serve such as specific resources & data sources.
- Harder to test and debug as a devleoper.
- Harder to maintain scalability is hard between container when it should be easier with prog specific feature

## What determine size of a cache ?
- Usage its being used and algorithm beng implemented to access and store data on it. 
- Hotness of request.
- How data lifecycling method is being handled.
- Cost of the RAM.
- The Working Set: A cache should ideally be sized based on the working set (the amount of unique data needed for your application to run smoothly at any given time). If your cache is smaller than your working set, you will experience "thrashing"—where the cache constantly evicts data that it needs to reload almost immediately.


# TCP and HTTP
### 1. TCP: The Reliable Transport Layer (The "Pipe")
TCP is strictly responsible for moving a **stream of bytes** from Point A to Point B reliably. It does not know that you are sending a "request" or a "video."

* **Connection Lifecycle:** TCP handles the **3-way handshake** (establishing the connection), **packet sequencing** (reordering chunks that arrive out of order), and **acknowledgments** (ensuring data wasn't lost).
* **Addressing:** It relies on the **IP layer** for machine addressing and uses **Ports** to ensure data reaches the correct software process, not the OS itself.
* **Limitations:** TCP is "application-agnostic." It treats all data as an undifferentiated stream. It does not know where one message ends and the next begins.



---

### 2. HTTP: The Application Layer (The "Language")
HTTP is the protocol that gives meaning to the bytes flowing through the TCP pipe. It sits *on top* of TCP. HTTP is what defines the "request" and the "type of data."

* **Message Framing:** Since TCP is just a stream of bytes, HTTP adds headers like `Content-Length` or `Transfer-Encoding: chunked`. This tells the receiver: *"Start here and stop after X bytes."* This is how the receiver "builds the chunks back out."
* **Metadata (Headers):** Things like `User-ID`, `Content-Type`, and `Authorization` are not handled by TCP; they are **Application Headers**. HTTP uses these to tell the server how to interpret the payload.
* **Semantics:** HTTP defines the **verbs** (GET, POST, PUT, DELETE). TCP has no concept of a "GET" request; it just sees characters like `G`, `E`, `T`.



---

### Key Improvements to Your Mental Model

| Concept | Your Original Thought | Technical Correction |
| :--- | :--- | :--- |
| **TCP's Role** | "Hardware layer," "Handles request types" | TCP is the **Transport Layer**. It is agnostic to request types. |
| **Data Types** | "What is the incoming data type" | This is defined by the **HTTP `Content-Type` header**, not TCP. |
| **Chunking** | "How to build chunks back out" | TCP chops data for **network efficiency**; HTTP frames data for **logical meaning**. |
| **The "Pipe"** | TCP as "Electricity" | TCP is the **delivery truck**; the bytes are the **cargo**; HTTP is the **manifest** inside the box. |

### Revised Summary
You can think of it as a hierarchy:
1.  **IP Layer:** Finds the building (the machine).
2.  **TCP Layer:** Finds the office room (the port) and ensures the delivery person doesn't drop the boxes (reliability).
3.  **HTTP Layer:** Opens the box, reads the manifest (headers), and identifies what the package is (video, JSON, HTML) so the application knows what to do with it.


## OOP Domination
- Easier to understand to mmimic the business-case that it reuqire to use
- Ecnourage extreme code reusability and function calling
- Most of the app being created are open to extension and changes based on biz requirement, a procedural or event typically will produced better performance and sizes but it requires massive refactoring too change, even if we decided to add feature we will sandbox it in oop manner.
- Writing and maintaining what you write is easier to do in oop.

## Problem with Go
- Weird syntax ofr public or private funcs
- Opinionated async/distributed work (making it bad for ml optimization)
- inability to do a procedural transformation or default output of error type.
- Lack of good libs for other than tui, dsa, or backend work.

## Why we do Inversion of Control ?
- So each object can be enough decoupled such that we can test it separately or do a object simulation.
- There is no need for the code to instantiate an object over and over everytime a method is called.
- most ioc framework has an opinionated way for dfiferent inheritor to sync up with each other making development of feature and synchronization easier.


Moving from a legacy MySQL setup to a modern, event-driven PostgreSQL architecture is a high-stakes operation. To pull this off without data loss or extended downtime, you need a coordinated "dance" between the old and the new.

Here is the technical step-by-step execution plan for your migration:

### Phase 1: The "Audit & Schema" (Preparation)
Before touching a single row of data, you must align the structures.

1.  **Schema Conversion:** Use a tool (like `pgLoader` or `db-forge`) to map MySQL types to Postgres. 
    * *Watch out:* Convert MySQL `TINYINT(1)` to Postgres `BOOLEAN` and `DATETIME` to `TIMESTAMPTZ`.
2.  **Constraint Mapping:** Manually review Foreign Keys and Indexes. Postgres is stricter about data integrity than MySQL; if your MySQL data has "orphaned" rows, the Postgres migration will fail.
3.  **Dependency Mapping:** Identify every API, cron job, and reporting tool that touches the MySQL DB. This is your "Impact Map."

### Phase 2: The "Bridge" (Application Changes)
This is where you implement the **Dual-Write** or **Outbox Pattern** to ensure no new data is lost during the transition.

4.  **Abstraction Layer:** Refactor your DB interaction code to use the **Repository Pattern** or **IoC** (as we discussed). This allows your app to talk to "A Database" without caring if it's MySQL or Postgres.
5.  **Implement Message Queues (Kafka/RabbitMQ):** * Set up a "CDC" (Change Data Capture) tool like **Debezium**. 
    * As changes happen in MySQL, Debezium streams them to Kafka. 
    * A "Consumer" function reads from Kafka and writes that same data into Postgres.
    
6.  **Dual-Read Validation:** (Optional but recommended) Have your API read from *both* databases for 1% of traffic and compare the results to ensure the logic is identical.

### Phase 3: The "Deep Sync" (Historical Data)
Now you move the "Old Data" that existed before you turned on the Message Queue.

7.  **Snapshot Migration:** Take a "point-in-time" dump of the MySQL data. 
8.  **Transform & Load:** Use an ETL process to clean the data (fixing those `0000-00-00` dates) and load it into Postgres.
9.  **Catch-up:** Once the historical data is in, your Kafka/RabbitMQ stream (from Phase 2) will "replay" all the changes that happened during the snapshot, bringing Postgres up to the exact millisecond of MySQL.

### Phase 4: The "Testing & Stress"
10. **Unit & Integration Tests:** Run your entire test suite against the Postgres instance.
11. **Stress Testing:** Use a tool like `pgbench` or replay your MySQL slow-query logs against Postgres. 
    * *Goal:* Ensure Postgres's Query Planner is handling your complex joins as efficiently as MySQL did.
    

### Phase 5: The "Canary & Cutover"
12. **The "Quiet Period":** Based on your analytics, pick the window with the lowest traffic.
13. **Canary Deployment:** Update 5% of your API instances to use Postgres as the **Primary** (Source of Truth). Monitor error rates and latency closely.
14. **Final Cutover:** Route 100% of traffic to the Postgres-connected APIs. 
15. **Decommission:** Keep the MySQL instance in "Read-Only" mode for 48 hours as a safety net, then shut it down.

### Phase 6: Optimization (The "Extra Features")
16. **Post-Migration Cleanup:** Now that you are on Postgres, implement those advanced features:
    * Use **JSONB** for semi-structured data.
    * Implement **PostGIS** if you have location data.
    * Fine-tune your **Vacuum** settings and autovacuum triggers for long-term health.

---

### Pro-Tip: The "Rollback" Trigger
Before you start Phase 5, define a "Point of No Return." For example: *"If latency increases by >200ms or error rates hit 1% for more than 5 minutes, we flip the traffic back to the MySQL-only API branch."*

**Which of these phases feels like the biggest "bottleneck" for your current team?** (Usually, it's either the Schema Mapping or the Stress Testing).
