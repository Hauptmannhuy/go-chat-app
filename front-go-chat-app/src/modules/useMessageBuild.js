import { useState } from "react"

export function useMessageBuild(){

  const [messages, setMessageStorages] = useState({})

  function addMessageStorage(chatName){
    setMessageStorages((messages) => {
      const newMessages = {...messages}
      newMessages[chatName] = []
      return newMessages
    })
  }



  function addMessage(message){
    setMessageStorages((prevMessages) => {
      const newMessages = {...prevMessages}
      newMessages[message.chat_name].push(message)
      return newMessages
    })
  }

  function handleMessageLoad(data) {
    const chats = Object.keys(data)

    chats.forEach((chatName) => {
      const chatMessages = data[chatName]
      chatMessages.forEach(message => {
        addMessage(message)
      });
    })
    
  }


  return {messages, addMessageStorage, addMessage, handleMessageLoad}
}