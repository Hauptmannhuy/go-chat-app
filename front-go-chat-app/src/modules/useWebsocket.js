import { useEffect, useRef } from "react";

export function useWebsocket(url, onMessage) {
  const socket = useRef(null)

  useEffect(() => {
    const websocket = new WebSocket(url)

    websocket.addEventListener("open", () => {
      socket.current = websocket
    })

    websocket.addEventListener("message", onMessage)

    return () => websocket.close()
  },  [])
  
   const sendMessage = (message) => {
    socket.current.send(JSON.stringify(message))
  }

  return { sendMessage }
}