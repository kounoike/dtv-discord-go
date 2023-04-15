console.log(`hello, world`)

import mysql, { RowDataPacket } from "mysql2/promise"
import Fuse from "fuse.js"
import fs from "fs"

interface Program {
  id: number
  name: string
  description: string
  startAt: Date
  channel: string
  extended: string
}

interface IProgram extends RowDataPacket {
  id: number
  name: string
  description: string
  start_at: number
  channel: string
  message_id: string
  json: string
}

async function main() {
  const connection = await mysql.createConnection({
    host: "db",
    user: "dtv-discord",
    password: "dtv-discord",
    database: "dtv",
  })

  const [rows, fields] = await connection.query<IProgram[]>({
    sql: `
    SELECT
      program.*,
      service.name AS channel,
      program_message.message_id
    FROM program
      JOIN service ON program.service_id = service.service_id
        AND program.network_id = service.service_id
      JOIN program_message ON program.id = program_message.program_id
      `,
  })

  const parseExtended = (extended: Object) => {
    return JSON.stringify(extended)
  }

  const document = rows.map((row): Program => {
    return {
      id: row.id,
      name: row.name.normalize(),
      description: row.description.normalize(),
      channel: row.channel.normalize(),
      startAt: new Date(row.start_at),
      extended: parseExtended(JSON.parse(row.json).extended).normalize(),
    }
  })

  const options: Fuse.IFuseOptions<Program> = {
    keys: ["name", "description", "json"],
    includeScore: true,
    includeMatches: true,
    shouldSort: true,
  }

  const index = Fuse.createIndex(options.keys ?? [], document)
  fs.writeFileSync("document.json", JSON.stringify(document))
  fs.writeFileSync("index.json", JSON.stringify(index.toJSON()))

  const fuse = new Fuse<Program>(document, options, index)

  console.log(fuse.search("トニカク").map((e) => [e.item.name, e.score]))

  connection.destroy()
}

await main()
