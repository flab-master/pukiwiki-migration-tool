import { useState } from 'react'
import './App.css'

type FormData = {
  userID: string
  pageID: string
}

function App() {
  const [formData, setFormData] = useState<FormData>({
    userID: "",
    pageID: "",
  })

  const addData = (e: React.ChangeEvent<HTMLInputElement>) => {
    const key = e.target.name as keyof FormData
    const value = e.target.value

    setFormData((current) => ({
      ...current,
      [key]: value,
    }))
  }

  const transApply = async (userData: FormData) => {
    console.log(userData.userID + "が移行申請しました。")
    if(!userData.userID.trim() || !userData.pageID.trim()){
      console.error("userIDかpageIDの両方を入力してください")
      return
    }

    try {
      const response = await fetch("http://localhost:8080/api/migrate", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          user: userData.userID,
        }),
      })

      if(!response.ok){
        console.error("移行申請に失敗しました")
        return
      }

      console.log(userData.userID + "が移行申請しました。")
    } catch (error){
      console.error("通信エラーが発生しました", error)
    }
  }

  return (
    <>
      <h1>Pukiwiki-Migration-Tool</h1>
      <h2>pukiwikiの個人のページのseminar-personal/userIDのuserIDの部分を入力</h2>
      <p>（例：fujimoto2025）</p>
      <input
       name="userID"
       type="text" 
       value={formData.userID}
       onChange={addData}
       placeholder='pukiwikiのuserIDを入力' 
      />

      <h2>ページを追加するnotionのページIDを入力</h2>
      <input
       name="pageID"
       type="text" 
       value={formData.pageID}
       onChange={addData}
       placeholder='notionのpageIDを入力' 
      />
      <br></br>

      <button onClick={() => transApply(formData)}>移行申請</button>
      <p>{formData.userID}</p>
      <p>{formData.pageID}</p>
    </>
  )
}

export default App
