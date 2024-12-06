
import { useState } from "react";



export function useChatBuild() {

  const [chats, setChats] = useState({});

  function addChat(name,id,participation, type)  {
    setChats((chats) => {
      let chatProperties = null
      if (type == 'group') {
        chatProperties = createNewChatObject(name,id,participation, type)
      } else {
        const displayName = name
        chatProperties = createNewChatObject(displayName, id, participation, type)
      }
      const newChats = {...chats}
      newChats[name] = chatProperties
      return newChats
    })
  }

  function handleNewGroupChat(chat) {
    const name = chat.chat_name
    const creator_id = chat.creator_id
    addChat(name, chat.chat_id, true, 'group')
  }

  const handleInitChatLoad = (chats, type,  participation = false) => {
    if (!chats) return 
    typeof chats != 'object' ? chats = [chats] : null
    const chatKeys = Object.keys(chats)

    chatKeys.forEach(chatName => {
      const chat = chats[chatName]
      addChat(chat.chat_name, chat.chat_id, participation, type)
     
    })
  }



  
  const createNewChatObject = (chatName, id ,participation, type) => ({name: chatName, id: id,participation: participation, type: type})


  return {chats, addChat, createNewChatObject, handleInitChatLoad, handleNewGroupChat }
}