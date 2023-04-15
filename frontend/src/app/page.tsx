"use client"
import Image from "next/image"
import { Inter } from "next/font/google"
import styles from "./page.module.css"
import {
  AppBar,
  IconButton,
  ListItem,
  ListItemText,
  TextField,
  Toolbar,
  Typography,
  Link,
  Box,
} from "@mui/material"
import { useAsync, useDebounce } from "react-use"
import Fuse from "fuse.js"
import { useEffect, useRef, useState } from "react"
import { FixedSizeList, ListChildComponentProps } from "react-window"

const inter = Inter({ subsets: ["latin"] })

const documentUrl = "/index/document.json"
const indexUrl = "/index/index.json"

interface Program {
  id: number
  name: string
  description: string
  json: string
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
  const resultRef = useRef<HTMLDivElement>(null)
  const [resultHeight, setResultHeight] = useState<number>(0)

  useEffect(() => {
    if (resultRef.current) {
      setResultHeight(resultRef.current.getBoundingClientRect().height)
    }
  }, [resultRef, query])

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

  const renderRow = ({ index, style }: ListChildComponentProps) => {
    if (results === undefined) return <></>
    return (
      <ListItem
        style={style}
        key={index}
        component="div"
        disablePadding
        className={styles.listItem}
      >
        <ListItemText
          primary={results[index].item.name}
          secondary={results[index].item.description}
        />
        <Link href="/">Discord Link</Link>
      </ListItem>
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
              <Typography variant="h6" className={styles.title}>
                視聴ちゃん {resultHeight}
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
              noValidate
              autoComplete="off"
              onSubmit={(e) => e.preventDefault()}
              className={styles.form}
            >
              <TextField
                id="search-program"
                label="番組検索"
                value={query}
                style={{ width: "100%" }}
                onChange={(e) => setQuery(e.target.value.normalize("NFKC"))}
              />
            </form>
            <div ref={resultRef} className={styles.result}>
              {results && (
                <FixedSizeList
                  height={resultHeight}
                  width="100%"
                  itemSize={46}
                  itemCount={results.length}
                  overscanCount={5}
                >
                  {renderRow}
                </FixedSizeList>
              )}
            </div>
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
