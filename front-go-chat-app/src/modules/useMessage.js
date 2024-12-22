import { useState } from "react"

export function useMessage(){

  const [messages, setMessageStorages] = useState({})

  function addMessageStorage(chatName){
    setMessageStorages((messages) => {
      const newMessages = {...messages}
      newMessages[chatName] = []
      return newMessages
    })
  }

  function addMessage(message){
    setMessageStorages((messages) => {
      const newMessages = {...messages}
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
 * @param {Array} data 
 */

  function handleMessageLoad(data) {
    const names = data.map((el) => (el.chat_name)).filter((val, index, self) => (index == self.indexOf(val,0)))

    names.forEach((chatName) => {
      addMessageStorage(chatName)
    })
    
    data.forEach(message => {
      addMessage(message)
    });
  }


  return {messages, addMessageStorage, addMessage, handleMessageLoad}
}