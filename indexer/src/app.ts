console.log(`hello, world`)

import mysql, { RowDataPacket } from "mysql2/promise"
import Fuse from "fuse.js"
import fs from "fs"
import { Logger, ILogObj } from "tslog"

const logger: Logger<ILogObj> = new Logger<ILogObj>({
  type: "pretty",
  minLevel: 0,
})

interface Program {
  id: number
  name: string
  description: string
  startAt: Date
  channel: string
  extended: string
}

interface Recording {
  programId: number
  name: string
  description: string
  startAt: Date
  channel: string
  extended: string
  m2tsPath: string
  mp4Path: string
  aribb24SubtitlePath: string
  transcribedPath: string
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

interface IRecordedFiles extends RowDataPacket {
  id: number
  program_id: number
  content_path: string
}

interface IJobEncoded extends RowDataPacket {
  id: number
  program_id: number
  output_path: string
}

interface IJobTranscribed extends RowDataPacket {
  id: number
  programId: number
  m2tsPath: string
  mp4Path: string
  aribB24SubtitlePath: string
  transcribedTextPath: string
}

async function initializeProgramDocument(connection: mysql.Connection) {
  const [rows, fields] = await connection.query<IProgram[]>({
    sql: `
    SELECT
      program.*,
      service.name AS channel,
      program_message.message_id
    FROM program
      JOIN service ON program.service_id = service.service_id
        AND program.network_id = service.network_id
      JOIN program_message ON program.id = program_message.program_id
      `,
  })

  const parseExtended = (extended?: Object) => {
    if (extended === undefined) return ""
    let ret = ""
    for (const [k, v] of Object.entries(extended)) {
      ret += `${k.normalize("NFKC")}: ${v.normalize("NFKC")}\n`
    }
    return ret
  }

  const document = rows.map((row): Program => {
    return {
      id: row.id,
      name: row.name.normalize("NFKC"),
      description: row.description.normalize("NFKC"),
      channel: row.channel.normalize("NFKC"),
      startAt: new Date(row.start_at),
      extended: parseExtended(JSON.parse(row.json).extended),
    }
  })

  const options: Fuse.IFuseOptions<Program> = {
    keys: ["name", "description", "json"],
    includeScore: true,
    includeMatches: true,
    shouldSort: true,
  }

  const index = Fuse.createIndex(options.keys ?? [], document)
  fs.writeFileSync(
    "/document_index/program_document.json",
    JSON.stringify(document)
  )
  fs.writeFileSync(
    "/document_index/program_index.json",
    JSON.stringify(index.toJSON())
  )

  const fuse = new Fuse<Program>(document, options, index)

  console.log(fuse.search("NEWS").map((e) => [e.item.name, e.score]))
}

async function initializeRecordedFiles() {}

async function main() {
  const connection = await mysql.createConnection({
    host: "db",
    user: "dtv-discord",
    password: "dtv-discord",
    database: "dtv",
  })

  await initializeProgramDocument(connection)

  connection.destroy()
}

await main()
