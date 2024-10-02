
import './App.css'
import ChatRoom from './components/ChatRoom';
import Home from './components/Home';
import { createBrowserRouter, RouterProvider, Link } from 'react-router-dom'
function App() {
  const router = createBrowserRouter([
    {
    path: '/',
    element: <Home/>,
  },
  {
    path: '/home',
    element: <Element/>
  },
  {
    path: '/chat',
    element: <ChatRoom/>
  }
])
  return (
    <>
    
    <RouterProvider router={router}/>
    
    </>
  )
}

export default App
