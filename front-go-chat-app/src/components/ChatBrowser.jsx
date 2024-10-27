import { useEffect, useState } from "react";

import Chat from "./Chat";
import ChatList from "./ChatList";
import { useNavigate } from "react-router-dom";
import SignOutButton from "./SignOutButton";
import { useRef } from "react";
import Search from "./Search";

function ChatBrowser(){ 
   

  const socketConnection = useRef(null)

  const navigate = useNavigate()
  
  const [chats, setChats] = useState([]);
  const [messages, setMessages] = useState({})


  const [chatSelected, setSelectedChat] = useState(null)
  const [creationChatInvoked, setInvokeStatus] = useState(false)
  const [creationInput, setCreationInput] = useState("")

  const [searchResults, setSearchResults] = useState(null)
  
  const [lastResponse, setLastResponse] = useState(null)


  const sendEnvelope = (type, data) => {
      const actions = {
        "NEW_MESSAGE": {
          type: "NEW_MESSAGE",
          user_id: `${data[0]}`,
          chat_id: `${data[1]}`,
          body: `${data[2]}`,
        },
        "NEW_CHAT": {
          type: "NEW_CHAT",
          chat_id: `${data[0]}`,
          user_id: `${data[1]}`
        },
        "JOIN_CHAT": {
          type: "JOIN_CHAT",
          chat_id: `${data[0]}`,
          user_id: `${data[1]}`
        },
        "SEARCH_QUERY": {
          type: "SEARCH_QUERY",
          input: `${data[0]}`
        }
      }
    
      const json = actions[type];
      if (json) {
        socketConnection.current.send(JSON.stringify(json));
      } else {
        console.error(`Unknown action type: ${type}`);
      }
    };

  const createChat = (name,type = "NEW_CHAT", user_id = getUsernameCookie()) => {
    setInvokeStatus(false)
    setCreationInput(false)
    
    sendEnvelope(type, [name, user_id])
  };

  const joinChat = (name,type = "JOIN_CHAT") => {
    appendNewChat(name)
    sendEnvelope(type, [name, getUsernameCookie() ])
  }
  
  const appendNewChat = (name) => {
    const previuosChats = [...chats]
    const newChats = previuosChats.filter(el => (el.name != name))
    const changedChat = {name: name, participation: true}
    newChats.push(changedChat)
    setChats(newChats)
  }

  const addChatsAndMessages = (names, participation = false) => {
    if (!names) return
    typeof names != 'object' ? names = [names] : null
    names.forEach(name => {
      addChatHandler(name, participation)
      addMessagesObjectHandler(name)
    })
  }

  
  const addChatHandler = (name, participation) => {
    setChats((chats) => {
      console.log(chats)
      const newChat = createNewChatObject(name, participation)
      return [...chats, newChat]
    })
  }

  const createNewChatObject = (name, participation) => ({name: `${name}`, participation: participation})

  const addMessagesObjectHandler = (name) => {
    setMessages((messages) => {
      const newMessages = {...messages}
      newMessages[name] = []
      return newMessages
    })
  }

  function saveLocalMessage(message) {
    setMessages((messages) => {
      const newMessages = {...messages}
      console.log(newMessages)
      
    
      newMessages[message.chat_id].push(message)
      return newMessages
    })
  }

  function handleMessageLoad(data){
   setMessages((messages) => ({...messages, ...data}))
  }


  const sendMessage = (chatID, userID, msg, type = "NEW_MESSAGE") => {
    sendEnvelope(type, [userID, chatID, msg])
  }
  
 
  function processSocketMessage(ev) {
    const response = JSON.parse(ev.data);
    console.log(response.Data)
    switch (response.Type) {
      case "NEW_CHAT":
        addChatsAndMessages(response.Data.chat_id, true);
        break;
      case "NEW_MESSAGE":
        saveLocalMessage(response.Data)
        break;
      case "LOAD_SUBS":
        return addChatsAndMessages(response.Data, true);
      case "LOAD_MESSAGES":
        handleMessageLoad(response.Data)
        console.log(response)
        break;
      case "SEARCH_QUERY":
        console.log(response.Data)
        setLastResponse(response.Data)
        break;
      case "ERROR":
        console.log(response)
      default:
        console.log(ev.data.Data);
        break;
    }
  }; 

 

  const search = (input) => {
    if (input == ""){
      setLastResponse(null)
      return
    }
    sendEnvelope("SEARCH_QUERY", [input])
  }

  function handleSearchQuery(data){
    
    const filterChats = (data) => {
      if (!data.chats) return []

      const participatedChats = chats.filter(chat => chat.participation);

      const participatedQueriedChats = participatedChats.filter((el) =>  data.chats.includes(el.name));
      const participatingChatNames = participatedQueriedChats.map(chat => chat.name);

       const newChats = data.chats
        .filter(chatName => !participatingChatNames.includes(chatName))
        .map(name => createNewChatObject(name, false));    

      return [...newChats, ...participatedQueriedChats]
    }
    const fetchedChats = filterChats(data)
    console.log(fetchedChats)
    setSearchResults(fetchedChats)
  }
  
  useEffect(() => {
    if (lastResponse){
      handleSearchQuery(lastResponse)
    }
  }, [lastResponse])

  const userAuthenticated = () => {
    if (document.cookie != '') return false;
    return true;
  }

  const getUsernameCookie = () =>  document.cookie.split('=')[1]

  
 
  useEffect(()=>{
    const websocket = new WebSocket("/socket/chat")

    websocket.addEventListener("open", () => {
      socketConnection.current = websocket
    })

    websocket.addEventListener("message", processSocketMessage)

    
    websocket.addEventListener("close", (ev) => {
     
    })
    return () => websocket.close()
  }, [])



  return (
    <>

    {userAuthenticated ? (< SignOutButton/>) : null}

    <div className="search-section">

    <p>
    <Search
    searchHandler={search}/>
    </p>

   { searchResults && lastResponse ?  
    <ChatList 
     chats={searchResults}
     handleSelect={setSelectedChat}
     handleJoin={joinChat}
     /> : 
     <ChatList
      chats={chats}
      handleSelect={setSelectedChat}
      handleJoin={joinChat}
      />
   } 
    </div>

    <div className="chat-bar">

    

    <button onClick={() => setInvokeStatus(true)}>Create chat</button>
    
    {creationChatInvoked ? (
      <div>
        <input type="text" onChange={(e) => setCreationInput(e.target.value)}/>
        <button onClick={() => createChat(creationInput)}>Create</button>
      </div>
    ) : (
      null
    )}
    </div>
    
    <div className="chat-display">
      {chatSelected ? (
      <Chat
      chatName = {chatSelected}
      msgHandler = {sendMessage}
      messages = {messages[chatSelected]}
      userID = {getUsernameCookie()}/>
      ) : (
      <h2>Chat display</h2>
          )}
    </div>
     
  </>
  )
}

export default ChatBrowser