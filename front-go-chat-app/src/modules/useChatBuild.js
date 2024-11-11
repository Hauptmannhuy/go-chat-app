
import { useState } from "react";



export function useChatBuild(addMessageStorage) {

  const [chats, setChats] = useState({});

  function addChat(name,id,participation, type)  {
    setChats((chats) => {
      console.log(name, participation, type)
      let chatProperties = null
      if (type == 'group') {
        console.log(name)
        chatProperties = createNewChatObject(name,id,participation, type)
      } else {
        const displayName = name
        chatProperties = createNewChatObject(displayName, id ,participation, type)
      }
      const newChats = {...chats}
      newChats[name] = chatProperties
      return newChats
    })
    addMessageStorage(name)
  }

  function handleNewGroupChat(chat) {
    console.log(chat)
    const name = chat.chat_name
    const creator_id = chat.creator_id
    addChat(name, chat.chat_id, true, 'group')
  }

  const handleInitChatLoad = (chats, type,  participation = false) => {
    console.log(chats)
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