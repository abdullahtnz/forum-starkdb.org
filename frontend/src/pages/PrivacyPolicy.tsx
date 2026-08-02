export default function PrivacyPolicy() {
  return (
    <div className="min-h-screen pt-20 pb-12">
      <div className="max-w-3xl mx-auto px-4 sm:px-8">
        <div className="card">
          <div className="section-tag">Legal</div>
          <h1 className="text-3xl font-display font-extrabold text-stark-text tracking-[-0.03em] mb-8">Privacy Policy</h1>

          <div className="font-mono text-sm text-stark-muted leading-relaxed space-y-6">
            <p>Last updated: August 2026</p>

            <div>
              <h2 className="text-stark-text font-display font-bold text-lg mb-2">1. Information We Collect</h2>
              <p>We collect information you provide directly to us when you create an account, including your username, email address, and any content you post on the forum (questions, answers, comments). We also collect usage data such as IP addresses, browser type, and pages visited for analytics and security purposes.</p>
            </div>

            <div>
              <h2 className="text-stark-text font-display font-bold text-lg mb-2">2. How We Use Your Information</h2>
              <p>Your information is used to provide and maintain the STARKDB Forum service, including account management, content display, and community moderation. We use your email address for account-related communications only. We do not send marketing emails without your explicit consent.</p>
            </div>

            <div>
              <h2 className="text-stark-text font-display font-bold text-lg mb-2">3. Data Storage and Security</h2>
              <p>Your data is stored on secure servers using industry-standard encryption. We implement appropriate technical and organizational measures to protect your personal information against unauthorized access, alteration, disclosure, or destruction. Passwords are hashed using bcrypt and never stored in plain text.</p>
            </div>

            <div>
              <h2 className="text-stark-text font-display font-bold text-lg mb-2">4. Cookies</h2>
              <p>STARKDB Forum uses essential cookies for authentication and session management. We use JWT tokens stored in your browser's local storage for maintaining login sessions. We do not use third-party tracking cookies.</p>
            </div>

            <div>
              <h2 className="text-stark-text font-display font-bold text-lg mb-2">5. Data Sharing</h2>
              <p>We do not sell, trade, or rent your personal information to third parties. Your published content (posts, comments) is publicly visible on the forum. Your email address is never displayed publicly.</p>
            </div>

            <div>
              <h2 className="text-stark-text font-display font-bold text-lg mb-2">6. Your Rights</h2>
              <p>You have the right to access, correct, or delete your personal data at any time through your profile settings. You may request complete account deletion by contacting the administrator. You can update your username from your profile settings.</p>
            </div>

            <div>
              <h2 className="text-stark-text font-display font-bold text-lg mb-2">7. Data Retention</h2>
              <p>We retain your account information and content for as long as your account is active. Deleted posts and comments are permanently removed. Upon account deletion, all associated data is removed from our systems within 30 days.</p>
            </div>

            <div>
              <h2 className="text-stark-text font-display font-bold text-lg mb-2">8. Contact</h2>
              <p>For privacy-related inquiries, please contact the STARKDB team. This forum is operated as part of the STARKDB open-source project by Abdullah Novruzlu.</p>
            </div>

            <div>
              <h2 className="text-stark-text font-display font-bold text-lg mb-2">9. Changes to This Policy</h2>
              <p>We may update this privacy policy from time to time. We will notify users of any material changes by posting a notice on the forum. Continued use of the forum after changes constitutes acceptance of the updated policy.</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
