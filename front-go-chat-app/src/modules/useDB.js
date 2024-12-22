import { useEffect, useRef } from "react";



export function useDB() {
  const indexDB = useRef(null)
  const dbName = 'test-db'
    
  async function connectDB(){
  return new Promise((resolve, reject) => {
    const openRequest = window.indexedDB.open(dbName)

    openRequest.addEventListener("success", () => {
      indexDB.current = openRequest.result
      console.log("connected to DB")
      console.log(indexDB.current)
      resolve('connect')
    })

      openRequest.addEventListener("error", () => {
      reject("error connect to DB")
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
      resolve('upgrade')
        })
    })      
  }
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

 async function cacheMessages(data){
    console.log("db:", indexDB)
    const messageStore = initDBtransaction("messages")
    console.log(data)
    const objectChatNames = Object.keys(data)
    for (let i = 0; i < objectChatNames.length; i++) {
      const chatName = objectChatNames[i]
      const chatObj = data[chatName]
      for (let j = 0; j < chatObj.length; j ++) {
        const newMessageObject = {message_id: chatObj[j].message_id, body: chatObj[j].body, chat_name: chatObj[j].chat_name, username: chatObj[j].username }
        const addRequest = messageStore.add(newMessageObject)

        addRequest.addEventListener("error", () => {
          console.log("Error adding new item to DB")
          console.log(addRequest.error)
        })
      }
    }
  }

  
 async function cacheChats(data){
    console.log(data)
    const groupChatsStore = initDBtransaction("group_chats")
    const privateChatsStore = initDBtransaction("private_chats")

    const cache = (data, storage) => {
      const keys = Object.keys(data)
      keys.forEach(name => {
        storage.add(data[name], data[name].chat_id)
      });
    }
    try {
      if (!data || !data.private){
        throw new Error("Private chat data is missing");
      }
      cache(data.private, privateChatsStore)
    } catch (error) {
      console.error(error)
    }

    try {
      if (!data || !data.private){
        cache(data.group, groupChatsStore)
        throw new Error("Group chat data is missing");
      }
    } catch (error) {
      console.error(error)
    }
   
  }

  function getChats(){
    console.log("get chats")
    const groupChatsStore = initDBtransaction("group_chats")
    const privateChatsStore = initDBtransaction("private_chats")
    return { privateChatReq: privateChatsStore.getAll(), groupChatReq: groupChatsStore.getAll() }
  }


  function getMessages() {
    const messageStore = initDBtransaction("messages")
    const request = messageStore.getAll()
    return request
  }
  
  function saveMessage(message){
    console.log("db mesg", message)
    const messageStore = initDBtransaction("messages")
   const req = messageStore.add(message)
   req.addEventListener("success", () =>{
    console.log("req success",req.result)
  })
  req.addEventListener("error",()=>{
    console.log("req success",req.error)
  })
  }

  function savePrivateChat(chat){
    console.log(chat)
    const privateChatsStore = initDBtransaction("private_chats")
    let req = privateChatsStore.add(chat, chat.chat_id)
    req.addEventListener("success", () =>{
      console.log("req success",req.result)
    })
    req.addEventListener("error",()=>{
      console.log("req success",req.error)
    })
  }




  return {connectDB, cacheMessages, cacheMessages, cacheChats, getMessages, getChats, saveMessage, savePrivateChat}
}