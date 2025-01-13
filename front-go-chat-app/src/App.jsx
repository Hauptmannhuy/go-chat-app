
import { useEffect, useState } from 'react';
import './App.css';
import ChatLayout from './components/ChatLayout';

import { GlobalContext } from './contexts/GlobalContext';

import { actionDispatcher } from './modules/actionDispatcher';
import { useChat } from './modules/useChat';
import { useDB } from './modules/useDB';
import { useMessage } from './modules/useMessage';
import { useSearchQuery } from './modules/useSearchQuery';
import { useWebsocket } from './modules/useWebsocket';
import useOnline from './modules/useOnline';

function App() {

  const {onlineStatus, changeOnlineStatus} = useOnline()
  const {searchResults, searchProfileResults, handleSearchQuery} = useSearchQuery()
  const {messages, addMessage, addMessageStorage, handleMessageLoad, initMessageStorages} = useMessage()
  const {chats, addChat, handleInitChatLoad } = useChat()
  const [selectedChat, selectChat] = useState(null)

  const {connectDB, saveMessage, savePrivateChat, saveGroupChat, cacheChats, cacheMessages, getChats, getMessages} = useDB(fetchCache)

  const {processSocketMessage, checkFetchStatus} = actionDispatcher({
    chatService: {add: addChat, initialLoad:handleInitChatLoad},
    messageService: {addMessage:addMessage, addStorage: addMessageStorage},
    searchService: {handleSearch: handleSearchQuery},
    dbService: { saveMessage: saveMessage, 
                 savePrivateChat: savePrivateChat, 
                 cacheChats: cacheChats, 
                 cacheMessages: cacheMessages, 
                 saveGroupChat: saveGroupChat, 
                 getChats: getChats, 
                 getMessages: getMessages
    },
    uiManager: {selectChat: selectChat},
    userService: {changeOnlineStatus: changeOnlineStatus}
  }, )

  const {writeToSocket, connectWS} = useWebsocket("/socket/chat", processSocketMessage)

 async function fetchCache() {
    console.log("asking cache")
    writeToSocket({type: "LOAD_SUBS"})
    writeToSocket({type: "LOAD_MESSAGES"})
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
      handleInitChatLoad(result, "private", true)
      initMessageStorages(result)
    })

    privateChatReq.addEventListener("error", () => {
      console.error("Error during indexDB request",privateChatReq.error)
    })

    groupChatReq.addEventListener("success", () => {
      const result = groupChatReq.result
      console.log("group chat result",result)
      handleInitChatLoad(result, "group", true)
      initMessageStorages(result)
    })

    groupChatReq.addEventListener("error", () => {
      console.error("Error during indexDB request",groupChatReq.error)

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
    <GlobalContext.Provider value={{writeToSocket, selectChat, selectedChat, messages, chats, searchProfileResults, searchResults, onlineStatus}}>
        <ChatLayout/>
    </GlobalContext.Provider>
    </>
  );
}

export default App;
