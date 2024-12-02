import { useEffect, useRef } from "react";

      /**
       * Returns objectStore.
       * @param {function(): Promise} onUpgrade 
       * @param {function} onConnect
       * @returns {void}
       */

export function useIndexedDB(onUpgrade, onConnect) {
  const indexDB = useRef(null)
  const dbName = 'test-db'
  useEffect(() => {

   async function connectDB(){

    console.log(onUpgrade, onConnect)
     const openRequest = window.indexedDB.open(dbName)
     openRequest.addEventListener("error", () => {
       console.log("error connect")
      })
      
      openRequest.addEventListener("success", () => {

        indexDB.current = openRequest.result
        // setIndexDB(openRequest.result)
        console.log("connected to IndexedDB")
          onConnect()
      })

      
      openRequest.addEventListener("upgradeneeded", () => {
        const db = openRequest.result
        console.log("upgrading")

        const messageStore = db.createObjectStore("messages",{keyPath: "id", autoIncrement: true}) 
        messageStore.createIndex("message_id", "message_id", {unique: false})
        messageStore.createIndex("body", "body", {unique: false})
        messageStore.createIndex("chat_name", "chat_name", {unique: false})
        messageStore.createIndex("user_id", "user_id", {unique: false} )

        const privateChatStore = db.createObjectStore("private_chats") 
        privateChatStore.createIndex("chat_id", "id", {unique: true})
        privateChatStore.createIndex("user1_id", "user1_id", {unique: false})
        privateChatStore.createIndex("user2_id", "user2_id", {unique: false})

        const groupChatStore = db.createObjectStore("group_chats") 
        groupChatStore.createIndex("chat_id", "id", {unique: true})
        groupChatStore.createIndex("name", "name", {unique: false})
        groupChatStore.createIndex("creator_id", "creator_id", {unique:false})
        console.log("need cache")
        
          onUpgrade()
          
      })

    }

    connectDB()
  }, [])

      /**
       * Returns objectStore.
       * @param {string} dbname
       * @returns {IDBObjectStore}
       */

   const initDBtransaction = (dbname) => { 
    const transaction = indexDB.current.transaction([dbname], "readwrite")
    const objectStore = transaction.objectStore(dbname)
    return objectStore
  }

  function cacheMessages(data){
    console.log("db:", indexDB)
    const messageStore = initDBtransaction("messages")
    console.log(data)
    // const parsedData = JSON.parse(data)
    const objectChatNames = Object.keys(data)
    for (let i = 0; i < objectChatNames.length; i++) {
      const chatName = objectChatNames[i]
      const chatObj = data[chatName]
      for (let j = 0; j < chatObj.length; j ++) {
        console.log(chatObj[j])
        const newMessageObject = {message_id: chatObj[j].message_id, body: chatObj[j].body, chat_name: chatObj[j].chat_name, username: chatObj[j].username }
        const addRequest = messageStore.add(newMessageObject)

        addRequest.addEventListener("error", () => {
          console.log("Error adding new item to DB")
          console.log(addRequest.error)
        })
      }
    }
  }

  
  function cacheChats(data){
    const groupChatsStore = initDBtransaction("group_chats")
    const privateChatsStore = initDBtransaction("private_chats")
    const privateKeys = Object.keys(data.private)
    const groupKeys= Object.keys(data.group)
    console.log(data)
    privateKeys.forEach(key => {
      const chat = data.private[key]
      groupChatsStore.add(chat, chat.chat_id)
    });
    groupKeys.forEach(key => {
      const chat = data.group[key]
      privateChatsStore.add(chat, chat.chat_id)
    });
  }

  function getChats(){
    console.log("get chats")
    const groupChatsStore = initDBtransaction("group_chats")
    const privateChatsStore = initDBtransaction("private_chats")
    return { privateChatReq: groupChatsStore.getAll(), groupChatReq: privateChatsStore.getAll() }
  }


  function getMessages() {
    const messageStore = initDBtransaction("messages")
    const request = messageStore.getAll()
    return request
  }
  
  function saveMessage(message){
    const messageStore = initDBtransaction("messages")
    messageStore.add(message)
  }




  return {indexDB,cacheMessages, cacheChats, getMessages, getChats, saveMessage}
}