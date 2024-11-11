
export function useChatAction(sendMessage, { handleInitChatLoad, addChat,handleMessageLoad, addMessage, handleSearchQuery}){

  

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
    console.log(ev.data)
    const actionOnType = {
      NEW_CHAT: () => {
        const {chat_name,chat_id } = response.Data
       addChat(chat_name,chat_id, 'group', true)
      },
      NEW_MESSAGE: () => {
        addMessage(response.Data)
      },
      NEW_PRIVATE_CHAT: () => {
        console.log(response)
        const {chat_name, chat_id, message} = response.Data
        addChat(chat_name,chat_id, 'private', true)
        addMessage(message)
      },
      LOAD_SUBS: () => {
        handleInitChatLoad(response.Data.group,'group', true)
        handleInitChatLoad(response.Data.private, 'private', true)
      },
       LOAD_MESSAGES: () => {
        handleMessageLoad(response.Data)
         console.log(response)
       },
       SEARCH_QUERY: () => {
         console.log(response.Data)
         handleSearchQuery(response.Data.SearchResults)
       },
       ERROR: () => {
         console.log(response)
       }
    }
    actionOnType[response.Type]()
  }
 
  
  
  



  return {sendEnvelope, processSocketMessage}
}

