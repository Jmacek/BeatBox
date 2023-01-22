import React, { useState, useEffect, useRef } from 'react';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemText from '@mui/material/ListItemText';
import { IconButton, ListItemButton } from '@mui/material';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import PauseIcon from '@mui/icons-material/Pause';
import { styled } from '@mui/material/styles';

import { apiUrl } from "../../configs/constants";
import { Song } from "../../configs/structures";

const Div = styled('div')(({ theme }) => ({
  ...theme.typography.button,
  backgroundColor: theme.palette.background.paper,
  padding: theme.spacing(1),
}));

interface Props {
    songs: Array<Song>;
}

const SongList: React.FC<Props> = ({ songs }) => {
  const [songId, setSongID] = useState<string | null>(null);

  const [isPlaying, setIsPlaying] = useState(false);

  const audioElement = useRef<HTMLAudioElement>(null);

  const fetchAndPlaySong = async () => {
      if (audioElement.current && songId != null) {
          audioElement.current.src = `${apiUrl}/api/music/stream?id=${songId}`;
          await audioElement.current.play();
      }
  };

  const handleStop = () => {
      if (audioElement.current) {
          audioElement.current.pause();
      }
  };

  const handleStart = async (id: string) => {
      // if we're resuming the same song simply play
    if (id != null && id == songId) {
        audioElement.current?.play();
    }
  else {
      handleStop();
      setSongID(id);
  }
  };

  useEffect(() => {
      fetchAndPlaySong();
      }, [songId]);

  useEffect(() => {
      const handlePlayPause = () => {
          setIsPlaying(!audioElement.current?.paused);
      };

      if (audioElement.current) {
          audioElement.current.addEventListener('play', handlePlayPause);
          audioElement.current.addEventListener('pause', handlePlayPause);
      }

    return () => {
          if (audioElement.current) {
              audioElement.current.removeEventListener('play', handlePlayPause);
              audioElement.current.removeEventListener('pause', handlePlayPause);
          }
      };
      }, [isPlaying, audioElement]);

  return (
      <>
      {songs ? (
          <List disablePadding>
            {songs.map((song) => (
              <ListItem
                key={song.SongID}
                disablePadding
                disableGutters
                >
                {isPlaying && song.SongID ==songId ? (
                        <ListItemButton divider onClick={() => handleStop()}>
                          <IconButton>
                            <PauseIcon />
                          </IconButton>
                          <ListItemText primary={`${song.ArtistName} - ${song.SongName} - ${song.AlbumName}`} />
                        </ListItemButton>
                        ) : (
                          <ListItemButton divider onClick={() => handleStart(song.SongID)}>
                            <IconButton>
                              <PlayArrowIcon />
                            </IconButton>
                            <ListItemText primary={`${song.ArtistName} - ${song.SongName} - ${song.AlbumName}`} />
                          </ListItemButton>
                          )}
              </ListItem>
              ))}
          </List>
          ) : (
            <Div>{"No Songs Found."}</Div>
          )}
      <audio ref={audioElement}/>
      </>
      )};


export default SongList;