
import { useEffect, useState } from 'react';
import './App.css';
import ChatLayout from './components/ChatLayout';
import SearchSection from './components/SearchSection';

import { GlobalContext } from './contexts/GlobalContext';

import { actionDispatcher } from './modules/actionDispatcher';
import { useChat } from './modules/useChat';
import { useDB } from './modules/useDB';
import { useMessage } from './modules/useMessage';
import { useSearchQuery } from './modules/useSearchQuery';
import { useWebsocket } from './modules/useWebsocket';

function App() {
  const {searchResults, searchProfileResults, handleSearchQuery} = useSearchQuery()
  const {messages, addMessage, addMessageStorage, handleMessageLoad} = useMessage()
  const {chats, addChat, handleInitChatLoad, createNewChatObject, handleNewGroupChat } = useChat()
  const [selectedChat, selectChat] = useState(null)

  const {connectDB, saveMessage, savePrivateChat, cacheChats, cacheMessages, getChats, getMessages} = useDB(fetchCache)

  const {processSocketMessage, checkFetchStatus} = actionDispatcher({
    chatService: {add: addChat, initialLoad:handleInitChatLoad},
    messageService: {addMessage:addMessage, addStorage: addMessageStorage},
    searchService: {handleSearch: handleSearchQuery},
    dbService: {saveMessage: saveMessage, savePrivateChat: savePrivateChat, cacheChats: cacheChats, cacheMessages: cacheMessages,},
    uiManager: {selectChat: selectChat}
  }, )

  const {sendMessage, connectWS} = useWebsocket("/socket/chat", processSocketMessage)


  

 async function fetchCache() {
    console.log("asking cache")
    sendMessage({type: "LOAD_SUBS"})
    sendMessage({type: "LOAD_MESSAGES"})
    checkFetchStatus()
    .then(() => {
      display()
    })
    .catch((reason) => {
      console.error("Error during cache fetch:", reason)
    })
  }

  async function display() {
    const {privateChatReq, groupChatReq} = getChats()
    const messageReq = getMessages()

    privateChatReq.addEventListener("success", () => {
      const result = privateChatReq.result
      console.log(privateChatReq)
      handleInitChatLoad(result, "private", true)
    })

    groupChatReq.addEventListener("success", () => {
      const result = groupChatReq.result
      handleInitChatLoad(result, "group", true)
    })
    
    messageReq.addEventListener("success", () => {
      handleMessageLoad(messageReq.result)
    })

}

    useEffect(() => {
      async function initializeApp() {
        let dbStatus = null
        try {
          dbStatus = await connectDB()
          if (dbStatus == 'connect') {
           await display()
          }
        } catch (error) {
          throw new Error("Error connecting to DB", error);
        }
        try {
         await connectWS()
         if (dbStatus == 'upgrade') await fetchCache()
        } catch(error) {
          throw new Error("Error connecting to WS:", error); 
        }
       }  
       initializeApp()
    },[])



  return (
    <>
    <GlobalContext.Provider value={{sendMessage, selectChat, selectedChat, messages, chats, searchProfileResults, searchResults}}>
        <ChatLayout/>
    </GlobalContext.Provider>
    </>
  );
}

export default App;
