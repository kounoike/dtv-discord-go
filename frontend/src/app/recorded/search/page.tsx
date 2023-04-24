"use client"
import Image from "next/image"
import { Inter } from "next/font/google"
import styles from "./page.module.css"
import {
  AppBar,
  IconButton,
  TextField,
  Toolbar,
  Typography,
  Link,
  ListItemIcon,
  MenuItem,
  Card,
  CardContent,
  CardActions,
  Button,
} from "@mui/material"
import { useAsync, useDebounce } from "react-use"
import Fuse from "fuse.js"
import { useEffect, useRef, useState } from "react"
import { Virtuoso } from "react-virtuoso"

const inter = Inter({ subsets: ["latin"] })

const documentUrl = "/index/recorded_document.json"
const indexUrl = "/index/recorded_index.json"

interface RecordedFiles {
  id: number
  programId: number
  m2tsPath?: string | null
  mp4Path?: string | null
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

const fuseOptions: Fuse.IFuseOptions<RecordedFiles> = {
  keys: [
    "name",
    "description",
    "extended",
    "genre",
    "channel",
    "aribB24Subtitle",
    "transcribedText",
  ],
  includeScore: true,
  includeMatches: true,
  shouldSort: true,
}

export default function SearchPage() {
  const asyncState = useAsync(async () => {
    const documentPromise = (async () => {
      const response = await fetch(documentUrl)
      const result = JSON.parse(await response.text()) as RecordedFiles[]
      return result
    })()
    const indexPromise = (async () => {
      const response = await fetch(indexUrl)
      return Fuse.parseIndex<RecordedFiles>(JSON.parse(await response.text()))
    })()
    const [document, index] = await Promise.all([documentPromise, indexPromise])

    const fuse = new Fuse(document, fuseOptions, index)
    return { document, fuse }
  })

  const [query, setQuery] = useState<string>("")
  const [debouncedQuery, setDebouncedQuery] = useState<string>("")
  const [results, setResults] = useState<Fuse.FuseResult<RecordedFiles>[]>()
  const resultRef = useRef<HTMLDivElement>(null)
  const [resultHeight, setResultHeight] = useState<number>(0)
  const defaultRef = useRef<HTMLDivElement>(null)
  const [defaultHeight, setDefaultHeight] = useState<number>(0)
  const [_, cancel] = useDebounce(
    () => {
      setDebouncedQuery(query)
    },
    300,
    [query]
  )

  useEffect(() => {
    if (resultRef.current) {
      setResultHeight(resultRef.current.getBoundingClientRect().height)
    }
  }, [resultRef, debouncedQuery])

  useEffect(() => {
    if (defaultRef.current) {
      setDefaultHeight(defaultRef.current.getBoundingClientRect().height)
    }
  }, [defaultRef, asyncState.value])

  useEffect(() => {
    if (asyncState.value === undefined) return
    if (debouncedQuery === "") {
      setResults(undefined)
      return
    }
    setResults(asyncState.value.fuse.search(debouncedQuery.normalize("NFKC")))
  }, [debouncedQuery, asyncState.value])

  const renderRow = (documentIndex: number) => {
    if (!asyncState.value) return

    const idx = asyncState.value.document.length - documentIndex - 1

    return (
      <Card sx={{ marginBottom: "1rem" }}>
        <CardContent>
          <Typography variant="h6" component="div">
            {asyncState.value.document[idx].name}
          </Typography>
          <Typography sx={{ fontSize: 16 }}>
            {asyncState.value.document[idx].description}
          </Typography>
          <Typography variant="body2">
            {asyncState.value.document[idx].extended}
          </Typography>
        </CardContent>
        <CardActions>
          {asyncState.value.document[idx].mp4Path && (
            <Button
              size="small"
              href={
                "/recorded/mp4/" +
                encodeURIComponent(asyncState.value.document[idx].mp4Path!)
              }
            >
              MP4
            </Button>
          )}
        </CardActions>
      </Card>
    )
  }

  const renderSearchRow = (index: number) => {
    if (results === undefined) return <></>
    return (
      <Card sx={{ marginBottom: "1rem" }}>
        <CardContent>
          <Typography variant="h6" component="div">
            {results[index].item.name}
          </Typography>
          <Typography sx={{ fontSize: 16, marginBottom: "1rem" }}>
            {results[index].item.description}
          </Typography>
          {results[index].matches?.map((m) => buildMatchValue(m))}
        </CardContent>
        <CardActions>
          {results[index].item.mp4Path && (
            <Button
              size="small"
              href={
                "/recorded/mp4/" +
                encodeURIComponent(results[index].item.mp4Path!)
              }
            >
              MP4
            </Button>
          )}
        </CardActions>
      </Card>
    )
  }

  const buildMatchValue = (match: Fuse.FuseResultMatch) => {
    if (match.value === undefined) return <></>

    let ret = [] as JSX.Element[]
    let prev = 0
    match.indices.slice(0, 10).forEach((tpl) => {
      ret.push(
        <>
          {Math.max(prev, tpl[0] - 20) <= 0 ? "" : "…"}
          {match.value?.slice(Math.max(prev, tpl[0] - 20), tpl[0])}
          <b>{match.value?.slice(tpl[0], tpl[1] + 1)}</b>
          {match.value?.slice(tpl[1] + 1, tpl[1] + 20)}…
        </>
      )
      // prev = tpl[1] + 1
    })

    const keyDescriptions: any = {
      name: "番組名",
      description: "番組説明",
      extended: "番組詳細",
      genre: "ジャンル",
      channel: "チャンネル",
      aribB24Subtitle: "字幕",
      transcribedText: "文字起こし",
    }

    return (
      <Typography variant="body2" color="text.primary2">
        <p>
          {match.key && <b key="title">{keyDescriptions[match.key]}: </b>}
          {ret}
        </p>
      </Typography>
    )
  }

  return (
    <>
      {asyncState.loading ? (
        <main className={styles.loadingMain}>
          <Image
            src="/shichou-chan-logo.png"
            alt="視聴ちゃんロゴ"
            className={styles.logo}
            width={1280}
            height={313}
            priority
          />
        </main>
      ) : (
        <>
          <AppBar position="sticky" className={styles.appbar}>
            <Toolbar>
              <Typography
                variant="h6"
                className={styles.title}
                sx={{ flexGrow: 1 }}
              >
                視聴ちゃん 録画検索 {asyncState.value?.document.length}件
              </Typography>
              <MenuItem
                onClick={() => (window.location.href = "/program/search")}
              >
                <Typography>→番組検索</Typography>
              </MenuItem>
              <MenuItem
                onClick={() => (window.location.href = "/recorded/search")}
              >
                <Typography>→録画検索</Typography>
              </MenuItem>
              <IconButton>
                <Image
                  alt="視聴ちゃんアイコン"
                  src="/web-icon.png"
                  width={50}
                  height={50}
                ></Image>
              </IconButton>
            </Toolbar>
          </AppBar>
          <main className={styles.main}>
            <form
              noValidate
              autoComplete="off"
              onSubmit={(e) => e.preventDefault()}
              className={styles.form}
            >
              <TextField
                id="search-program"
                label="録画検索"
                value={query}
                style={{ width: "100%" }}
                onChange={(e) => setQuery(e.target.value)}
              />
            </form>
            {results !== undefined && results.length > 0 ? (
              <div ref={resultRef} className={styles.result}>
                <Virtuoso
                  style={{ height: resultHeight }}
                  totalCount={results.length}
                  itemContent={renderSearchRow}
                />
              </div>
            ) : (
              <div ref={defaultRef} className={styles.result}>
                <Virtuoso
                  style={{ height: defaultHeight }}
                  totalCount={asyncState.value?.document.length ?? 0}
                  itemContent={renderRow}
                />
              </div>
            )}
          </main>
        </>
      )}
    </>
  )
}
