import { useContext, useRef } from "react";
import { GlobalContext } from "../contexts/GlobalContext";



export function ActionDispatcher({chatService, messageService, dbService, searchService}){

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
  

  

  const processSocketMessage = (ev) => {
    const response = JSON.parse(ev.data)
    console.log(response)
    const actionOnType = {
      NEW_CHAT: () => {
        const {chat_name,chat_id } = response.Data
        chatService.addChat(chat_name,chat_id, 'group', true)
      },
      NEW_MESSAGE: () => {
        dbService.saveMessage(response.Data)
        messageService.addMessage(response.Data)
      },
      NEW_PRIVATE_CHAT: () => {
        const {chat_name, chat_id, message, initiator_id} = response.Data
        chatService.add(chat_name,chat_id, true, 'private')
        messageService.addStorage(chat_name)
        messageService.addMessage({message: message, chat_name: chat_name, chat_id: chat_id})
        dbService.savePrivateChat({chat_name: chat_name, chat_id: chat_id})
        dbService.saveMessage({body: message, chat_name: chat_name, user_id: initiator_id, message_id: 0 }) 
      },
      LOAD_SUBS: async () => {
        await dbService.cacheChats(response.Data)
        fetchStatus.current.subStatus = true
      },
       LOAD_MESSAGES: async () => {
        await dbService.cacheMessages(response.Data)
        fetchStatus.current.messageStatus = true
        
       },
       SEARCH_QUERY: () => {
        console.log(searchService)
        searchService.handleSearch(response.Data.SearchResults)
       },
       ERROR: () => {
        return false
       }
    }
    actionOnType[response.Type]()
  }
 
 
  
  



  return { processSocketMessage, checkFetchStatus}
}

