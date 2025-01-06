import { useRef } from "react";



export function actionDispatcher({chatService, messageService, dbService, searchService, userService}){

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
  

  

  const processSocketMessage = async (ev) => {
    const response = JSON.parse(ev.data)
    console.log(response)
    const actionOnType = {
      NEW_GROUP_CHAT: async () => {
        const {chat_name,chat_id } = response.Data
        chatService.add(chat_name,chat_id, 'group', true)
        messageService.addStorage(chat_name)
        dbService.saveGroupChat(response.Data)
      },
      JOIN_CHAT: async () => {
        const {chat_name, chat_id, creator_id, message} = response.Data
        chatService.add(chat_name, chat_id, 'group', true)
        messageService.addStorage(chat_name)
        // messageService.addMessage(message)
        dbService.saveGroupChat({chat_name: chat_name, chat_id: chat_id, creator_id: creator_id})
        // dbService.saveMessage(message)
      },
      NEW_MESSAGE: async () => {
        dbService.saveMessage(response.Data)
        messageService.addMessage(response.Data)
      },
      NEW_PRIVATE_CHAT: async () => {
        const {chat_name, chat_id, body, initiator_id, init_username} = response.Data
        chatService.add(chat_name,chat_id, true, 'private')
        messageService.addStorage(chat_name)
        // messageService.addMessage({body: body, chat_name: chat_name, chat_id: chat_id, username: init_username, message_id: 0 })
        dbService.savePrivateChat({chat_name: chat_name, chat_id: chat_id})
        // dbService.saveMessage({body: body, chat_name: chat_name, user_id: initiator_id, message_id: 0 }) 
      },
      LOAD_SUBS: async () => {
        await dbService.cacheChats(response.Data)
        fetchStatus.current.subStatus = true
      },
       LOAD_MESSAGES: async () => {
        await dbService.cacheMessages(response.Data)
        fetchStatus.current.messageStatus = true
        
       },
       SEARCH_QUERY: async () => {
        searchService.handleSearch(response.Data.SearchResults)
        userService.changeOnlineStatus(response.Data.status)
       },
       OFFLINE_MESSAGES: async () => {
        const keys = Object.keys(response.Data)
        console.log(keys)
        keys.forEach((chatName) => {
          const messages = response.Data[chatName]
          for (const message of messages) {
            const {data, type} = message
            if (type == "NEW_MESSAGE") {
              messageService.addMessage(data)
              dbService.saveMessage(data)
            } else if (type == 'NEW_PRIVATE_CHAT') {
              const {body, chat_name, init_username, message_id, initiator_id, receiver_id, chat_id } = data
              chatService.add(data.chat_name,data.chat_id, true, 'private')
              messageService.addStorage(chat_name)
              messageService.addMessage({body: body, chat_name: chat_name, username: init_username, message_id: message_id})
              dbService.savePrivateChat({user1_id: initiator_id, user2_id: receiver_id, chat_id: chat_id, chat_name: chat_name})
              dbService.saveMessage({user_id: initiator_id, body: body, chat_name: chat_name, username: init_username}) 
            }
          }
        })
       },
       USER_STATUS: () => {
        userService.changeOnlineStatus(response.Data.Status)
       },
       ERROR: () => {
        fetchStatus.current.messageStatus = false
        return false
       }
    }
     await actionOnType[response.Type]()
  }
 
 
  
  



  return { processSocketMessage, checkFetchStatus}
}

