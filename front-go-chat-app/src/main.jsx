import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import App from './App.jsx'
import './index.css'
import AuthenticationForm from './components/AuthenticationForm.jsx'


const router = createBrowserRouter([
  {
    path:"/chats",
    element: <App/>
  },
  {
    path:"/",
    element: <App/>
  },
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
  <StrictMode>
    <RouterProvider router={router}/>
  </StrictMode>,
)
