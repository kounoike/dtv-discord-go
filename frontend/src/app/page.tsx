"use client"
import Image from "next/image"
import { Inter } from "next/font/google"
import styles from "./page.module.css"
import {
  AppBar,
  IconButton,
  InputBase,
  TextField,
  Toolbar,
  Typography,
} from "@mui/material"
import { useAsync, useDebounce } from "react-use"
import Fuse from "fuse.js"
import { stringify } from "querystring"
import { useEffect, useState } from "react"

const inter = Inter({ subsets: ["latin"] })

const documentUrl = "/document.json"
const indexUrl = "/index.json"

interface Program {
  id: number
  name: string
  description: string
}

const fuseOptions: Fuse.IFuseOptions<Program> = {
  keys: ["name", "description", "json"],
  includeScore: true,
  includeMatches: true,
  shouldSort: true,
}

export default function Home() {
  const asyncState = useAsync(async () => {
    const documentPromise = (async () => {
      const response = await fetch(documentUrl)
      const result = JSON.parse(await response.text()) as Program[]
      return result
    })()
    const indexPromise = (async () => {
      const response = await fetch(indexUrl)
      return Fuse.parseIndex<Program>(JSON.parse(await response.text()))
    })()
    const [document, index] = await Promise.all([documentPromise, indexPromise])

    const fuse = new Fuse(document, fuseOptions, index)
    return fuse
  })

  const [query, setQuery] = useState<string>("")
  const [debouncedQuery, setDebouncedQuery] = useState<string>("")
  const [results, setResults] = useState<Fuse.FuseResult<Program>[]>()

  const [, cancel] = useDebounce(
    () => {
      setDebouncedQuery(query)
    },
    250,
    [query]
  )

  useEffect(() => {
    if (asyncState.value === undefined) return
    if (debouncedQuery === "") {
      setResults(undefined)
      return
    }
    setResults(asyncState.value.search(debouncedQuery))
  }, [debouncedQuery, asyncState.value])

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
              <Typography variant="h6" className={styles.title}>
                視聴ちゃん
              </Typography>
              <div className={styles.grow} />
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
              style={{ width: "100%" }}
              noValidate
              autoComplete="off"
              onSubmit={(e) => e.preventDefault()}
            >
              <TextField
                id="search-program"
                label="番組検索"
                value={query}
                style={{ width: "100%" }}
                onChange={(e) => setQuery(e.target.value)}
              />
            </form>
            {results && (
              <ul>
                {results.slice(0, 20).map((result) => (
                  <li key={result.item.id}>
                    {result.item.name} score: {result.score ?? "???"} matches
                    len: {result.matches?.length}
                    <ul>
                      {result.matches?.map((match, idx) => (
                        <li key={result.item.id + "-" + idx}>
                          indices: {match.indices.join(",")} key: {match.key}{" "}
                          refIndex: {match.refIndex} value:{" "}
                          {buildMatchValue(match)}
                        </li>
                      ))}
                    </ul>
                  </li>
                ))}
              </ul>
            )}
          </main>
        </>
      )}
    </>
  )
}

function buildMatchValue(match: Fuse.FuseResultMatch) {
  if (match.value === undefined) return <></>

  let ret = Array<JSX.Element>()
  let prev = 0
  match.indices.forEach((tpl) => {
    ret.push(
      <>
        {match.value?.slice(prev, tpl[0])}
        <b>{match.value?.slice(tpl[0], tpl[1] + 1)}</b>
      </>
    )
    prev = tpl[1] + 1
  })

  return (
    <>
      {ret}
      {match.value?.slice(prev)}
    </>
  )
}
