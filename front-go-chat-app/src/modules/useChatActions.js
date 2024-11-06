export function useChatAction(sendMessage, {setChats, setMessages, setProfiles, setsearchProfileResults,setSearchResults}){


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
        addChatsAndMessages(response.Data, true)
      },
      NEW_MESSAGE: () => {
        saveLocalMessage(response.Data)
      },
      NEW_PRIVATE_CHAT: () => {
        console.log(response)
      },
      LOAD_SUBS: () => {
        addChatsAndMessages(response.Data.group,'group', true)
        addChatsAndMessages(response.Data.private, 'private', true)
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
    

  const addChatsAndMessages = (chats, type,  participation = false) => {
    console.log(chats)
    if (!chats) return 
    typeof chats != 'object' ? chats = [chats] : null
    const chatKeys = Object.keys(chats)

    chatKeys.forEach(chatName => {
      const chat = chats[chatName]
      addChatHandler(chat.chat_name, participation, type)
      addMessagesObjectHandler(chat.chat_name)
    })
  }
  
  const addChatHandler = (name, participation, type) => {
    setChats((chats) => {
      console.log("setting chats")
      let chatProperties = null
      if (type == 'group') {
        chatProperties = createNewGroupChatObject(name, participation, type)
      } else {
       const displayName = name.split('_').filter(el => (el != getUsernameCookie()))
       
        chatProperties = createNewPrivateChatObject(name, displayName[0], participation)
      }
      const newChats = {...chats}
      newChats[name] = chatProperties
      return newChats
    })
  }
  
    function saveLocalMessage(message) {
      setMessages((messages) => {
        const newMessages = {...messages}
        console.log(newMessages)
        console.log(newMessages)
        
        
        newMessages[message.chat_name].push(message)
        return newMessages
      })
    }

    const addMessagesObjectHandler = (name) => {
      setMessages((messages) => {
        const newMessages = {...messages}
        newMessages[name] = []
        return newMessages
      })
    }
    
    function handleMessageLoad(data){
      setMessages((messages) => ({...messages, ...data}))
    }

    function handleSearchQuery(data){
      console.log(data)
      const groupChats = queryGroupChats(data)
      const privateChats = queryProfiles(data)
      setSearchResults({...privateChats, ...groupChats})
      setsearchProfileResults(data.users)
    }

    const queryProfiles = (data) => {
      // if (data.users.length == 0) return []
      const newPrivateChats = {}
      const newProfiles = {}
      const keys = Object.keys(data.users)
       keys.forEach((key) => { 
        const userProfile = data.users[key].profile
        console.log(userProfile)
        newProfiles[key] = userProfile
       newPrivateChats[key] = createNewGroupChatObject(userProfile.username, false, 'private')
    })
      return newPrivateChats
    }
  
    const queryGroupChats = (data) => {
      if (!data.chats) return []
  
      // const participatedChats = chats.filter(chat => chat.participation);
  
      // const participatedQueriedChats = participatedChats.filter((el) =>  data.chats.includes(el.name));
      // const participatingChatNames = participatedQueriedChats.map(chat => chat.name);
  
       const newChats = {}
       data.chats.forEach(chat => newChats[chat.chat_id] = createNewGroupChatObject(chat.chat_id, false));    
  
      return {...newChats, ...chats}
    }
  
     
  const createNewGroupChatObject = (chatName, participation) => ({name: chatName, participation: participation, type: "group"})
  const createNewPrivateChatObject = (chatName, displayName, participation) => ({name: chatName, displayName: displayName, participation: participation, type: "private"})




  return {sendEnvelope, processSocketMessage}
}