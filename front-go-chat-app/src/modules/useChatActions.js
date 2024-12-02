import { useRef } from "react";



export function useChatAction(sendMessage, { cacheChats, addChat,cacheMessages, addMessage, handleSearchQuery, saveMessage}){
  const fetchStatus = useRef({messageStatus: null, subStatus: null })

  const delay = async (ms) => new Promise((resolve) => {setTimeout(resolve,ms) } )
  
  async function checkFetchStatus() {
    let times = 0
    while (!fetchStatus.current.messageStatus || !fetchStatus.current.subStatus){
      if (times > 100 ){
        throw new Error("Error fetching cache");
      } else {
        await delay(50)
      }
      times += 1
    
    }
    
    return true
  }
  

  const sendEnvelope = (type, data) => {
    const actions = {
      "NEW_MESSAGE": {
        type: "NEW_MESSAGE",
        chat_name: `${data[0]}`,
        body: `${data[1]}`,
      },
      "NEW_CHAT": {
        type: "NEW_CHAT",
        chat_name: `${data[0]}`,
      },
      "NEW_PRIVATE_CHAT": {
        type: "NEW_PRIVATE_CHAT",
        receiver_id: `${data[0]}`,
        message: `${data[1]}`,
      },
      "JOIN_CHAT": {
        type: "JOIN_CHAT",
        chat_id: `${data[0]}`,
      },
      "SEARCH_QUERY": {
        type: "SEARCH_QUERY",
        input: `${data[0]}`,
      }
    }
  
    const json = actions[type];
    if (json) {
      return sendMessage(json)
    } else {
      console.error(`Unknown action type: ${type}`);
    }
  };

  const processSocketMessage = (ev) => {
    const response = JSON.parse(ev.data)
    const actionOnType = {
      NEW_CHAT: () => {
        const {chat_name,chat_id } = response.Data
       addChat(chat_name,chat_id, 'group', true)
      },
      NEW_MESSAGE: () => {
        addMessage(response.Data)
        saveMessage(response.Data)
      },
      NEW_PRIVATE_CHAT: () => {
        const {chat_name, chat_id, message} = response.Data
        addChat(chat_name,chat_id, 'private', true)
        addMessage(message)
      },
      LOAD_SUBS: () => {
        console.log(response)
        cacheChats(response.Data)
        fetchStatus.current.subStatus = true
      },
       LOAD_MESSAGES: () => {
        cacheMessages(response.Data)
        fetchStatus.current.messageStatus = true
       },
       SEARCH_QUERY: () => {
         handleSearchQuery(response.Data.SearchResults)
       },
       ERROR: () => {
        return false
       }
    }
    actionOnType[response.Type]()
  }
 
  
  
  



  return {sendEnvelope, processSocketMessage, checkFetchStatus}
}

