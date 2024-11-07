import { useState } from "react"

export function useMessageBuild(){

  const [messages, setMessageStorages] = useState({})

  function addMessageStorage(chatName){
    setMessageStorages((messages) => {
      const newMessages = {...messages}
      newMessages[chatName] = []
      console.log(newMessages)
      return newMessages
    })
  }



  function addMessage(message){
    setMessageStorages((prevMessages) => {
      const newMessages = {...prevMessages}
      console.log(message, newMessages)
      newMessages[message.chat_name].push(message)
      return newMessages
    })
  }

  function handleMessageLoad(data) {
    const chats = Object.keys(data)
    console.log(chats)
    console.log(data)
    chats.forEach((chatName) => {
      const chatMessages = data[chatName]
      chatMessages.forEach(message => {
        addMessage(message)
      });
    })
    
  }


  return {messages, addMessageStorage, addMessage, handleMessageLoad}
}