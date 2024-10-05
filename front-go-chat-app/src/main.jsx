import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.jsx'
import Chat from './components/Chat.jsx'
import './index.css'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import Chats from './components/ChatBrowser.jsx'
import ChatBrowser from './components/ChatBrowser.jsx'

const router = createBrowserRouter([
  {
    path:"/",
    element: <App/>
  },
  {
    path:"/chats",
    element: <ChatBrowser/>
  },
  {
    path: "/chats/:chatID",
    element: <Chat/>
  }
])

createRoot(document.getElementById('root')).render(
  // <StrictMode>
    <RouterProvider router={router}/>
  // </StrictMode>,
)
