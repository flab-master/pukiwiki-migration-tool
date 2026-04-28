import { useState } from 'react'
import { useEffect } from 'react'
import './App.css'
import RequestPage from './components/requestPage'
import InputForm from './components/inputForm'
import {BrowserRouter, Routes, Route, Link, useNavigate} from 'react-router-dom'

function App() {

  //const navigate = useNavigate()

  return(
    <div>
      <BrowserRouter>
        <Link to='/input'>入力画面へ</Link>
        <Routes>
          <Route path='/' element={
            <>
              <h2>メイン画面</h2>
            </>
          } />
          <Route path='/input' element={<InputForm />} />
          <Route path='/request' element={<RequestPage />} />
        </Routes>

      </BrowserRouter>
    </div>
  )
}

export default App
