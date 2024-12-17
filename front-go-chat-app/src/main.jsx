import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.jsx'
import './index.css'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'

import ChatLayout from './components/ChatLayout.jsx'
import AuthenticationForm from './components/AuthenticationForm.jsx'

const router = createBrowserRouter([
  {
    path:"/chats",
    element: <App/>
  },
  // {
  //   path:"/chats",
  //   element: <ChatLayout/>
  // },
  {
    path: "/sign_up",
    element: <AuthenticationForm/>,
  },
  {
    path: "/sign_in",
    element: <AuthenticationForm/>,
  }
])

createRoot(document.getElementById('root')).render(
  // <StrictMode>
    <RouterProvider router={router}/>
  //  </StrictMode>,
)
