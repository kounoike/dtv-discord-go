"use client"
import styles from "./page.module.css"
import "instantsearch.css/themes/satellite.css"

import {
  AppBar,
  IconButton,
  Toolbar,
  Typography,
  MenuItem,
  Card,
  Grid,
  Button,
} from "@mui/material"
import { useEffect, useState } from "react"
import Image from "next/image"
import { instantMeiliSearch } from "@meilisearch/instant-meilisearch"
import {
  ClearRefinements,
  Configure,
  Highlight,
  InfiniteHits,
  InstantSearch,
  Pagination,
  Panel,
  RefinementList,
  SearchBox,
  Snippet,
  Stats,
} from "react-instantsearch-dom"

export default function Home() {
  const [searchClient, setSearchClient] = useState<Object | undefined>(
    undefined
  )
  useEffect(() => {
    const url = `${window.location.protocol}//${window.location.host}:7443`
    setSearchClient(
      instantMeiliSearch(url, undefined, {
        primaryKey: "id",
      })
    )
  }, [])

  return (
    <>
      <AppBar position="sticky" className={styles.appbar}>
        <Toolbar>
          <Typography
            variant="h6"
            className={styles.title}
            sx={{ flexGrow: 1 }}
          >
            視聴ちゃん 番組検索
          </Typography>
          <MenuItem onClick={() => (window.location.href = "/program/search")}>
            <Typography>→番組検索</Typography>
          </MenuItem>
          <MenuItem onClick={() => (window.location.href = "/recorded/search")}>
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
        {searchClient !== undefined && (
          <InstantSearch indexName="program" searchClient={searchClient}>
            <Configure
              queryLanguages={["ja"]}
              naturalLanguages={["ja"]}
              advancedSyntax={true}
            />
            <Grid container spacing={2}>
              <Grid item xs={3}>
                <div>
                  <Typography variant="h5">チャンネル</Typography>
                  <RefinementList attribute="チャンネル名" limit={100} />
                </div>
                <div>
                  <Typography variant="h5">ジャンル</Typography>
                  <RefinementList attribute="ジャンル" limit={100} />
                </div>
                <ClearRefinements />
              </Grid>
              <Grid item xs={9}>
                <SearchBox />
                <Stats />
                <InfiniteHits hitComponent={Hit} />
              </Grid>
            </Grid>
          </InstantSearch>
        )}
      </main>
    </>
  )
}

const toHourMinute = (time: number) => {
  const hour = Math.floor(time / 60 / 60 / 1000)
  const minute = Math.floor((time - hour * 60 * 60 * 1000) / 60 / 1000)
  if (hour > 0 && minute > 0) {
    return `${hour}時間${minute}分`
  } else if (hour > 0) {
    return `${hour}時間`
  } else {
    return `${minute}分`
  }
}

const Hit = ({ hit }: { hit: any }) => (
  <Panel>
    <Typography variant="h6">
      <Highlight attribute="タイトル" hit={hit} />
    </Typography>
    <Typography variant="subtitle1">
      <Snippet attribute="番組説明" hit={hit} />
    </Typography>
    <Typography variant="body1">
      <Snippet attribute="番組詳細" hit={hit} />
    </Typography>
    <div style={{ display: "flex" }}>
      <Typography>
        <b>放送局:</b> {hit.チャンネル名}{" "}
      </Typography>
      <Typography sx={{ marginLeft: "2rem" }}>
        <b>開始日時:</b>
        {new Date(hit.StartAt).toLocaleString()}
      </Typography>
      <Typography sx={{ marginLeft: "2rem", flexGrow: 1 }}>
        <b>放送時間:</b>
        {toHourMinute(hit.Duration)}
      </Typography>
      <Button variant="text" href={hit.DiscordMessageUrl}>
        Discord
      </Button>
      <Button variant="text" target="_blank" href={hit.WebMessageUrl}>
        Web
      </Button>
    </div>
  </Panel>
)
