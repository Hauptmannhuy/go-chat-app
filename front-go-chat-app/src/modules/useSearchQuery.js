import { useState } from "react"
import { useChat } from "./useChat"

import { useDB } from "./useDB";



export function useSearchQuery(){
  const { createNewChatObject } = useChat()
  const [searchResults, setSearchResults] = useState({})
  const [searchProfileResults, setSearchProfileResults] = useState({})

  


  function handleSearchQuery(data){
    if (!data || data.length === 0) return

    const newProfiles = {}
    const newProfileChats = {}
    const newGroupChats = {}

    const profileKeys = Object.keys(data.users)
    const groupChatKeys = Object.keys(data.chats)
    const groupChats = data.chats
    const users = data.users
    
    for (let i = 0; i < profileKeys.length; i++) {
      const profileName = profileKeys[i]
      const profile = users[profileName].profile
      const chatProfile = users[profileName].chat
      console.log(chatProfile)
      newProfiles[profileName] = profile
      if (chatProfile.handshake) {
        newProfileChats[chatProfile.chat_name] = createNewChatObject(
          chatProfile.chat_name,
          chatProfile.chat_id,
          true, 
          'private'
        )
      } else {
        newProfileChats[profileName] = createNewChatObject(
          profileName,
          profile.id,
          false, 
          'private'
        )
      }
     
    }
    
    for (let i = 0; i < groupChatKeys.length; i++) {
      const groupChatName = groupChatKeys[i];
      const groupChat = groupChats[groupChatName]
      newGroupChats[groupChat.chat_name] = createNewChatObject(
        groupChat.chat_name, 
        groupChat.chat_id, 
        groupChat.is_subscribed == 'true' ? true : false,
         'group'
      )
    }
    setSearchProfileResults(() => ({  ...newProfiles }))
    setSearchResults(() => ({  ...newProfileChats, ...newGroupChats }))
}

  return {searchResults, searchProfileResults, handleSearchQuery}
}