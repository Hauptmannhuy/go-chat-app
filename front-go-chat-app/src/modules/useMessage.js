import { useState } from "react"
import { useDB } from "./useDB"

export function useMessage(){
 

  const [messages, setMessageStorages] = useState({})

  function addMessageStorage(chatName) {
    setMessageStorages((prevMessages) => {
      const newMessages = { ...prevMessages }
      newMessages[chatName] = []
      return newMessages
    })
  }

  /**
   * 
   * @param {Message} message 
   */
  
  function addMessage(message){
    
    setMessageStorages((prevMessages) => {
      const newMessages = {...prevMessages}
      const chatMessages = newMessages[message.chat_name]
      const isDublicate = chatMessages.some((val) => (val.message_id == message.message_id))
      if (isDublicate) {
        return newMessages
      }

      newMessages[message.chat_name].push(message)
      return newMessages
    })
  }


 /**
 * @param {import("./useDB").Message[]} data 
 */

  function handleMessageLoad(data) {
    data.forEach(message => {
      addMessage(message)
    })
  }

   /**
   * @param {import("./useDB").Message[]} data 
   */

  function initMessageStorages(data) {
    console.log("initializing message storage", data)
    data.forEach((chat) => {
      addMessageStorage(chat.chat_name)
    })
  }



  return {messages, addMessageStorage, addMessage, handleMessageLoad, initMessageStorages}
}