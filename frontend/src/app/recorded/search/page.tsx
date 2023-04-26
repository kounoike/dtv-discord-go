"use client"
import Image from "next/image"
import styles from "./page.module.css"
import "instantsearch.css/themes/satellite.css"

import {
  AppBar,
  IconButton,
  Toolbar,
  Typography,
  MenuItem,
  Grid,
  Button,
} from "@mui/material"
import Accordion from "@mui/material/Accordion"
import AccordionSummary from "@mui/material/AccordionSummary"
import AccordionDetails from "@mui/material/AccordionDetails"
import { useEffect, useState } from "react"
import { instantMeiliSearch } from "@meilisearch/instant-meilisearch"
import {
  ClearRefinements,
  Configure,
  Highlight,
  InfiniteHits,
  InstantSearch,
  Pagination,
  RefinementList,
  SearchBox,
  Snippet,
  Stats,
} from "react-instantsearch-dom"
import { ExpandMore } from "@mui/icons-material"

export default function SearchPage() {
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
            視聴ちゃん 録画検索
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
          <InstantSearch indexName="recorded_file" searchClient={searchClient}>
            <Configure
              attributesToSnippet={["ARIB字幕:50", "文字起こし:50"]}
              queryLanguages={["ja"]}
              naturalLanguages={["ja"]}
              advancedSyntax={true}
            />
            <Grid container spacing={2}>
              <Grid item xs={3}>
                <div>
                  <Typography variant="h5">ジャンル</Typography>
                  <RefinementList
                    attribute="ジャンル"
                    style={{ listStyleType: "none" }}
                    limit={100}
                  />
                </div>
                <div>
                  <Typography variant="h5">チャンネル</Typography>
                  <RefinementList attribute="チャンネル名" />
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
  <div>
    <Typography variant="h6">
      <Highlight attribute="タイトル" hit={hit} />
    </Typography>
    <Typography variant="subtitle1">
      <b>番組説明:</b>
      <Snippet attribute="番組説明" hit={hit} />
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
      {hit.mp4 && (
        <Button
          variant="text"
          href={"/recorded/mp4/" + encodeURIComponent(hit.mp4)}
          target="_blank"
        >
          MP4
        </Button>
      )}
    </div>
    <Typography variant="body1">
      <Snippet attribute="番組詳細" hit={hit} />
    </Typography>
    {hit["ARIB字幕"] && (
      <Accordion>
        <AccordionSummary
          expandIcon={<ExpandMore />}
          id="ARIB-Subtitle"
          aria-controls="arib-subtitle"
        >
          <Typography>ARIB字幕</Typography>
        </AccordionSummary>
        <AccordionDetails>
          <Highlight attribute="ARIB字幕" hit={hit} />
        </AccordionDetails>
      </Accordion>
    )}
    {hit["文字起こし"] && (
      <Accordion>
        <AccordionSummary
          expandIcon={<ExpandMore />}
          id="Transcribed-Subtitle"
          aria-controls="transcibed-subtitle"
        >
          <Typography>文字起こし</Typography>
        </AccordionSummary>
        <AccordionDetails>
          <Snippet attribute="文字起こし" hit={hit} nbWords={20} />
        </AccordionDetails>
      </Accordion>
    )}
  </div>
)
