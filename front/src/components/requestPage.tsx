import { useEffect } from 'react'
import { useState } from 'react'
import './requestPage.css'
import { useLocation, useNavigate } from 'react-router-dom'

type Props = {
    userID: string
    pageID: string
}

export default function RequestPage(){
    const location = useLocation()
    const navigate = useNavigate()
    const props = location.state as Props

    const transApply = async (userData: Props) => {
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
                    pageID: userData.pageID
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

    return(
        <div className="reqContents">
            <h2>以下の内容で移行申請します。</h2>
            <p>userID: {props.userID}</p>
            <p>pageID: {props.pageID}</p>
            <button onClick={() => transApply(props)}>申請</button>
            <button onClick={() => navigate("/input")}>戻る</button>
        </div>
    )
}