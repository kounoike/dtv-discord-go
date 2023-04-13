console.log(`hello, world`)

import mysql, { RowDataPacket } from "mysql2/promise"
import Fuse from "fuse.js"
import fs from "fs"

interface Program extends RowDataPacket {
  id: number
  name: string
  description: string
}

async function main() {
  const connection = await mysql.createConnection({
    host: "db",
    user: "dtv-discord",
    password: "dtv-discord",
    database: "dtv",
  })

  const [rows, fields] = await connection.query<Program[]>({
    sql: "SELECT * FROM program",
  })

  const options: Fuse.IFuseOptions<Program> = {
    keys: ["name", "description", "json"],
    includeScore: true,
    includeMatches: true,
    shouldSort: true,
  }

  const index = Fuse.createIndex(options.keys ?? [], rows)
  fs.writeFileSync("document.json", JSON.stringify(rows))
  fs.writeFileSync("index.json", JSON.stringify(index.toJSON()))

  const fuse = new Fuse(rows, options, index)

  console.log(fuse.search("トニカク").map((e) => [e.item.name, e.score]))

  connection.destroy()
}

await main()
