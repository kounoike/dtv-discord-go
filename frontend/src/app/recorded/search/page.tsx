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
import { useEffect, useState } from "react"
import { instantMeiliSearch } from "@meilisearch/instant-meilisearch"
import {
  ClearRefinements,
  Highlight,
  Hits,
  InstantSearch,
  Pagination,
  RefinementList,
  SearchBox,
  Snippet,
  Stats,
} from "react-instantsearch-dom"

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
          <InstantSearch indexName="recorded_file" searchClient={searchClient}>
            <Grid container spacing={2}>
              <Grid item xs={3}>
                <div>
                  <Typography variant="h5">ジャンル</Typography>
                  <RefinementList
                    attribute="ジャンル"
                    style={{ listStyleType: "none" }}
                  />
                </div>
                <div>
                  <Typography variant="h5">チャンネル</Typography>
                  <RefinementList attribute="チャンネル名" />
                </div>
                <ClearRefinements />
              </Grid>
              <Grid item xs={9}>
                <Stats />
                <SearchBox />
                <Hits hitComponent={Hit} />
                <Pagination showLast={true} />
              </Grid>
            </Grid>
          </InstantSearch>
        )}
      </main>
    </>
  )
}

const Hit = ({ hit }: { hit: any }) => (
  <div>
    <Typography variant="h6">
      <Highlight attribute="タイトル" hit={hit} />
    </Typography>
    <Typography variant="subtitle1">
      <Snippet attribute="番組説明" hit={hit} />
    </Typography>
    <Typography variant="body1">
      <Snippet attribute="番組詳細" hit={hit} />
    </Typography>
    <div>
      <Button variant="text" href={"/recorded/mp4/" + hit.mp4} target="_blank">
        MP4
      </Button>
    </div>
  </div>
)
