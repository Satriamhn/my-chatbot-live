import { BrowserRouter as Router, Routes, Route } from "react-router";
import SignIn from "./pages/AuthPages/SignIn";
import SignUp from "./pages/AuthPages/SignUp";
import NotFound from "./pages/OtherPage/NotFound";
import UserProfiles from "./pages/UserProfiles";
import AppLayout from "./layout/AppLayout";
import { ScrollToTop } from "./components/common/ScrollToTop";
import Home from "./pages/Dashboard/Home";

// Chatbot UI pages
import BotSettings from "./pages/Settings/BotSettings";
import KnowledgeBase from "./pages/Knowledge/KnowledgeBase";
import Inbox from "./pages/Inbox/Inbox";
import Contact from "./pages/Contact/Contact";

export default function App() {
  return (
    <>
      <Router>
        <ScrollToTop />
        <Routes>
          {/* Dashboard Layout */}
          <Route element={<AppLayout />}>
            <Route index path="/" element={<Home />} />

            {/* Chatbot Pages */}
            <Route path="/settings/bot" element={<BotSettings />} />
            <Route path="/knowledge" element={<KnowledgeBase />} />
            <Route path="/inbox" element={<Inbox />} />
            <Route path="/contact" element={<Contact />} />

            {/* Others Page */}
            <Route path="/profile" element={<UserProfiles />} />
          </Route>

          {/* Auth Layout */}
          <Route path="/signin" element={<SignIn />} />
          <Route path="/signup" element={<SignUp />} />

          {/* Fallback Route */}
          <Route path="*" element={<NotFound />} />
        </Routes>
      </Router>
    </>
  );
}
