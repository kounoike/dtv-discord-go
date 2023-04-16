import mysql, { RowDataPacket, OkPacket } from "mysql2/promise"
import Fuse from "fuse.js"
import fs from "fs"
import { Logger, ILogObj } from "tslog"
import path from "path"
import { on } from "events"

const transcribedBasePath =
  process.env["TRANSCRIBED_BASE_PATH"] ?? "/transcribed"

const logger: Logger<ILogObj> = new Logger<ILogObj>({
  type: "pretty",
  minLevel: 0,
  prettyLogTimeZone: "local",
})

logger.info("indexer started.")

interface Program {
  id: number
  name: string
  description: string
  startAt: Date
  channel: string
  extended: string
}

interface RecordedFiles {
  id: number
  programId: number
  m2tsPath?: string
  mp4Path?: string
  aribB24SubtitlePath?: string
  transcribedTextPath?: string
  aribB24Subtitle?: string
  transcribedText?: string
  startAt: Date
  duration: number
  genre: string
  extended: string
  name: string
  description: string
  channel: string
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
  programId: number
  m2tsPath?: string
  mp4Path?: string
  aribB24SubtitlePath?: string
  transcribedTextPath?: string
  startAt: Date
  duration: number
  genre: string
  json: string
  name: string
  description: string
  channel: string
}

async function initializeProgramDocument(connection: mysql.Connection) {
  const [status] = await connection.query<IInvalidStatus[]>(
    `
      SELECT
        status
      FROM index_invalid
      where type = "program"`
  )
  if (status.length === 0) {
    const [ret, _] = await connection.execute<OkPacket>(
      `
        INSERT INTO index_invalid
          (type, status)
        VALUES
          ("program", "valid")
        ON DUPLICATE KEY UPDATE
          status = VALUES(status)`
    )
  }
  if (status.length === 1 && status[0].status !== "invalid") {
    return
  }

  logger.info("update program document and index started.")

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

  const [ret, _] = await connection.execute<OkPacket>(
    `
      INSERT INTO index_invalid
        (type, status)
      VALUES
        ("program", "valid")
      ON DUPLICATE KEY UPDATE
        status = VALUES(status)`
  )

  logger.info("update program document and index done.")
}

interface IInvalidStatus extends RowDataPacket {
  status: string
}

async function initializeRecordedFiles(connection: mysql.Connection) {
  const [status] = await connection.query<IInvalidStatus[]>(
    `
      SELECT
        status
      FROM index_invalid
      where type = "recorded"`
  )
  if (status.length === 0) {
    const [ret, _] = await connection.execute<OkPacket>(
      `
        INSERT INTO index_invalid
          (type, status)
        VALUES
          ("recorded", "valid")
        ON DUPLICATE KEY UPDATE
          status = VALUES(status)`
    )
  }
  if (status.length === 1 && status[0].status !== "invalid") {
    return
  }

  logger.info("update recorded document and index started.")

  const [rows, fields] = await connection.query<IRecordedFiles[]>({
    sql: `
    SELECT
      recorded_files.program_id as programId,
      recorded_files.m2ts_path as m2tsPath,
      recorded_files.mp4_path as mp4Path,
      recorded_files.aribb24_txt_path as aribB24SubtitlePath,
      recorded_files.transcribed_txt_path as transcribedTextPath,
      program.name,
      program.description,
      program.start_at,
      program.duration,
      program.json,
      program.genre,
      service.name AS channel
    FROM recorded_files
      JOIN program ON recorded_files.program_id = program.id
      JOIN service ON program.service_id = service.service_id AND program.network_id = service.network_id
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

  const document = rows.map((row): RecordedFiles => {
    const obj: RecordedFiles = {
      id: row.id,
      programId: row.programId,
      name: row.name.normalize("NFKC"),
      description: row.description.normalize("NFKC"),
      genre: row.genre,
      channel: row.channel.normalize("NFKC"),
      m2tsPath: row.m2tsPath,
      mp4Path: row.mp4Path,
      aribB24SubtitlePath: row.aribB24SubtitlePath,
      transcribedTextPath: row.transcribedTextPath,
      startAt: new Date(row.start_at),
      duration: row.duration,
      extended: parseExtended(JSON.parse(row.json).extended),
    }
    if (obj.aribB24SubtitlePath) {
      const f = fs.openSync(
        path.join(transcribedBasePath, obj.aribB24SubtitlePath),
        "r"
      )
      obj.aribB24Subtitle = fs.readFileSync(f, "utf-8")
      fs.closeSync(f)
    }
    if (obj.transcribedTextPath) {
      const f = fs.openSync(
        path.join(transcribedBasePath, obj.transcribedTextPath),
        "r"
      )
      obj.transcribedText = fs.readFileSync(f, "utf-8")
      fs.closeSync(f)
    }
    return obj
  })

  const options: Fuse.IFuseOptions<RecordedFiles> = {
    keys: [
      "name",
      "description",
      "genre",
      "channel",
      "extended",
      "aribB24Subtitle",
      "transcribedText",
    ],
    includeScore: true,
    includeMatches: true,
    shouldSort: true,
  }

  const index = Fuse.createIndex(options.keys ?? [], document)
  fs.writeFileSync(
    "/document_index/recorded_document.json",
    JSON.stringify(document)
  )
  fs.writeFileSync(
    "/document_index/recorded_index.json",
    JSON.stringify(index.toJSON())
  )

  const [ret, _] = await connection.execute<OkPacket>(
    `
      INSERT INTO index_invalid
        (type, status)
      VALUES
        ("recorded", "valid")
      ON DUPLICATE KEY UPDATE
        status = VALUES(status)`
  )
  // TODO: check

  logger.info("update recorded document and index done.")
}

async function main() {
  const connection = await mysql.createConnection({
    host: "db",
    user: "dtv-discord",
    password: "dtv-discord",
    database: "dtv",
  })

  await initializeProgramDocument(connection)
  await initializeRecordedFiles(connection)

  connection.destroy()
}

const cancelId = setInterval(main, 1000 * 60)

process.on("SIGINT", () => {
  clearInterval(cancelId)
  process.exit(0)
})
