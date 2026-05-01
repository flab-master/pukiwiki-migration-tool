import './App.css'
import RequestPage from './components/requestPage'
import InputForm from './components/inputForm'
import RequestListPage from './components/requestListPage'
import { Routes, Route, useNavigate } from 'react-router-dom'

function App() {
  const navigate = useNavigate()

  return (
    <div className="app-shell">
      <header className="app-header">
        <button className="ghost-button" onClick={() => navigate('/')}>
          ホーム
        </button>
        <button className="ghost-button" onClick={() => navigate('/requests')}>
          申請一覧
        </button>
        <button className="primary-button compact-button" onClick={() => navigate('/input')}>
          新規申請
        </button>
      </header>

      <Routes>
        <Route
          path="/"
          element={
            <main className="home-page">
              <div className="page-heading">
                <p className="eyebrow">Migration Tool</p>
                <h1>PukiWiki 移行申請</h1>
                <p>
                  PukiWiki の個人ページを Notion へ移行するための申請情報を入力します。
                </p>
              </div>

              <div className="home-actions">
                <button className="primary-button" onClick={() => navigate('/input')}>
                  入力を開始
                </button>
                <button className="secondary-button" onClick={() => navigate('/requests')}>
                  申請一覧を見る
                </button>
              </div>
            </main>
          }
        />
        <Route path="/input" element={<InputForm />} />
        <Route path="/request" element={<RequestPage />} />
        <Route path="/requests" element={<RequestListPage />} />
      </Routes>
    </div>
  )
}

export default App
