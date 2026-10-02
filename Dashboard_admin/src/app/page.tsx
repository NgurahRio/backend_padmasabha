"use client";
import { useEffect, useMemo, useRef, useState, type ReactNode } from "react";
import { api, idOf, unwrap } from "@/lib/api";
const nav = [
  ["overview", "Overview"],
  ["facility", "Facilities"],
  ["category", "Categories"],
  ["subcategory", "Sub Categories"],
  ["sos", "SOS"],
  ["user", "Users"],
];
const destinationNav = [
  ["villa", "Villas"],
  ["activity", "Activities"],
  ["restaurant", "Restaurants"],
  ["event", "Events"],
];
export default function Page() {
  const [ready, setReady] = useState(false),
    [user, setUser] = useState<any>(null),
    [active, setActive] = useState("overview"),
    [menu, setMenu] = useState(false),
    [destinationsOpen, setDestinationsOpen] = useState(true);
  useEffect(() => {
    const expired = () => setUser(null);
    window.addEventListener("admin-session-expired", expired);
    api("/admin/session")
      .then((d) => setUser(d.user))
      .catch(() => setUser(null))
      .finally(() => setReady(true));
    return () => window.removeEventListener("admin-session-expired", expired);
  }, []);
  if (!ready) return <div className="loader" />;
  if (!user) return <Login onLogin={setUser} />;
  return (
    <div className="shell">
      <aside className={menu ? "sidebar open" : "sidebar"}>
        <div className="brand">
          <img src="/padmanabha-logo.png" alt="" />
          <span>Padmanabha</span>
          <small>Manager</small>
        </div>
        <nav>
          {nav.slice(0, 1).map(([key, label]) => (
            <button
              key={key}
              className={active === key ? "selected" : ""}
              onClick={() => {
                setActive(key);
                setMenu(false);
              }}
            >
              <Icon name={key} />
              <span>{label}</span>
            </button>
          ))}
          <button className={destinationNav.some(([key]) => key === active) ? "nav-group-toggle selected" : "nav-group-toggle"} onClick={() => setDestinationsOpen((open) => !open)} aria-expanded={destinationsOpen}>
            <Icon name="destination" /><span>Destinations</span><Icon name="chevrondown" />
          </button>
          {destinationsOpen && <div className="sidebar-subnav">
            {destinationNav.map(([key, label]) => <button key={key} className={active === key ? "selected" : ""} onClick={() => { setActive(key); setMenu(false); }}><Icon name={key} /><span>{label}</span></button>)}
          </div>}
          {nav.slice(1).map(([key, label]) => (
            <button key={key} className={active === key ? "selected" : ""} onClick={() => { setActive(key); setMenu(false); }}><Icon name={key} /><span>{label}</span></button>
          ))}
        </nav>
        <button
          className="logout"
          onClick={async () => {
            try {
              await api("/admin/logout", { method: "POST" });
            } finally {
              setUser(null);
            }
          }}
        >
          <Icon name="logout" />
          <span>Logout</span>
        </button>
      </aside>
      {menu && <button className="backdrop" onClick={() => setMenu(false)} />}
      <button className="menu" onClick={() => setMenu((v) => !v)}>
        <Icon name="menu" />
      </button>
      <main>
        <Header user={user} />
        <AdminPage active={active} />
      </main>
    </div>
  );
}
function AdminPage({ active }: { active: string }) {
  const pages: Record<string, React.ReactNode> = {
    overview: <OverviewPage />,
    villa: <VillasPage />,
    activity: <ActivitiesPage />,
    restaurant: <RestaurantsPage />,
    event: <EventsPage />,
    category: <CategoriesPage />,
    subcategory: <SubCategoriesPage />,
    user: <UsersPage />,
    sos: <SOSPage />,
    facility: <FacilitiesPage />,
  };
  return pages[active] || <OverviewPage />;
}
function Login({ onLogin }: { onLogin: (u: any) => void }) {
  const [loginId, setLoginId] = useState(""),
    [password, setPassword] = useState(""),
    [rememberDevice, setRememberDevice] = useState(false),
    [busy, setBusy] = useState(false),
    [error, setError] = useState(""),
    [forgot, setForgot] = useState(false);
  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setBusy(true);
    setError("");
    try {
      const d = await api("/admin/login", {
        method: "POST",
        body: JSON.stringify({ loginId, password, rememberDevice }),
      });
      onLogin(d.user);
    } catch {
      setError("Email atau password salah!");
    } finally {
      setBusy(false);
    }
  }
  if (forgot) return <ForgotPassword onBack={() => setForgot(false)} />;
  return (
    <div className="login-page">
      <Decor />
      <form className="login-card" onSubmit={submit}>
        <Logo />
        <Heading
          tag="Welcome back"
          title="ADMIN LOGIN"
          text="Sign in to manage your Padmanabha dashboard"
        />
        <label>
          Email address
          <div className="input-shell">
            <MailIcon />
            <input
              type="email"
              placeholder="admin@example.com"
              autoComplete="email"
              value={loginId}
              onChange={(e) => setLoginId(e.target.value)}
              required
            />
          </div>
        </label>
        <label>
          Password
          <PasswordInput
            value={password}
            onChange={setPassword}
            placeholder="Enter your password"
          />
        </label>
        <div className="login-options">
          <label className="remember-check">
            <input
              type="checkbox"
              checked={rememberDevice}
              onChange={(e) => setRememberDevice(e.target.checked)}
            />
            <span />
            Remember this device
          </label>
          <button
            type="button"
            className="link-button forgot-link"
            onClick={() => setForgot(true)}
          >
            Forgot password?
          </button>
        </div>
        {error && <Error text={error} />}
        <Submit busy={busy} busyText="Signing in..." text="Login" />
        <p className="auth-footer">
          Secure access for Padmanabha administrators
        </p>
      </form>
    </div>
  );
}
function ForgotPassword({ onBack }: { onBack: () => void }) {
  const [step, setStep] = useState<1 | 2 | 3 | 4>(1),
    [email, setEmail] = useState(""),
    [otp, setOtp] = useState(""),
    [resetToken, setResetToken] = useState(""),
    [password, setPassword] = useState(""),
    [confirmPassword, setConfirmPassword] = useState(""),
    [busy, setBusy] = useState(false),
    [error, setError] = useState("");
  async function requestOTP(e?: React.FormEvent) {
    e?.preventDefault();
    setBusy(true);
    setError("");
    try {
      await api("/admin/forgot-password", {
        method: "POST",
        body: JSON.stringify({ email: email.trim().toLowerCase() }),
      });
      setStep(2);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setBusy(false);
    }
  }
  async function verify(e: React.FormEvent) {
    e.preventDefault();
    setBusy(true);
    setError("");
    try {
      const d = await api("/admin/verify-reset-otp", {
        method: "POST",
        body: JSON.stringify({ email, otp }),
      });
      setResetToken(d.resetToken);
      setStep(3);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setBusy(false);
    }
  }
  async function reset(e: React.FormEvent) {
    e.preventDefault();
    if (password !== confirmPassword) {
      setError("Konfirmasi password tidak sama");
      return;
    }
    setBusy(true);
    try {
      await api("/admin/reset-password", {
        method: "POST",
        body: JSON.stringify({ resetToken, password }),
      });
      setStep(4);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setBusy(false);
    }
  }
  return (
    <div className="login-page">
      <Decor />
      <div className="login-card reset-card">
        <Logo />
        <div className="step-dots">
          {[1, 2, 3].map((n) => (
            <i key={n} className={step >= n && step < 4 ? "active" : ""} />
          ))}
        </div>
        {step === 1 && (
          <form onSubmit={requestOTP}>
            <Heading
              tag="Account recovery"
              title="FORGOT PASSWORD"
              text="Masukkan email admin. Kami akan mengirim kode OTP yang berlaku selama 10 menit."
            />
            <label>
              Email address
              <div className="input-shell">
                <MailIcon />
                <input
                  type="email"
                  placeholder="admin@example.com"
                  autoComplete="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                />
              </div>
            </label>
            <Submit busy={busy} busyText="Sending OTP..." text="Send OTP" />
          </form>
        )}
        {step === 2 && (
          <form onSubmit={verify}>
            <Heading
              tag="Check your inbox"
              title="VERIFY OTP"
              text={`Kode OTP telah dikirim ke ${email}.`}
            />
            <label>
              6-digit OTP
              <OTPInput value={otp} onChange={setOtp} />
            </label>
            <Submit busy={busy} busyText="Verifying OTP..." text="Verify OTP" />
            <button
              type="button"
              className="link-button centered"
              onClick={() => requestOTP()}
            >
              Resend OTP
            </button>
          </form>
        )}
        {step === 3 && (
          <form onSubmit={reset}>
            <Heading
              tag="Almost finished"
              title="NEW PASSWORD"
              text="Gunakan minimal 8 karakter agar akun tetap aman."
            />
            <label>
              New password
              <PasswordInput
                value={password}
                onChange={setPassword}
                placeholder="Minimum 8 characters"
              />
            </label>
            <label>
              Confirm password
              <PasswordInput
                value={confirmPassword}
                onChange={setConfirmPassword}
                placeholder="Repeat new password"
              />
            </label>
            <Submit busy={busy} busyText="Saving..." text="Change Password" />
          </form>
        )}
        {error && <Error text={error} />}{" "}
        {step === 4 && (
          <div className="success-state">
            <span>
              <Icon name="check" />
            </span>
            <h1>PASSWORD CHANGED</h1>
            <p>Password berhasil diganti.</p>
            <button className="login-button" onClick={onBack}>
              Back to Login
            </button>
          </div>
        )}
        {step !== 4 && (
          <button className="link-button centered back-link" onClick={onBack}>
            <Icon name="arrowleft" /> Back to Login
          </button>
        )}
      </div>
    </div>
  );
}
function PasswordInput({
  value,
  onChange,
  placeholder,
}: {
  value: string;
  onChange: (v: string) => void;
  placeholder: string;
}) {
  const [visible, setVisible] = useState(false);
  return (
    <div className="input-shell">
      <LockIcon />
      <input
        type={visible ? "text" : "password"}
        value={value}
        placeholder={placeholder}
        minLength={8}
        onChange={(e) => onChange(e.target.value)}
        required
      />
      <button
        type="button"
        className="password-toggle"
        onClick={() => setVisible((v) => !v)}
      >
        <Eye off={visible} />
      </button>
    </div>
  );
}
function OTPInput({
  value,
  onChange,
}: {
  value: string;
  onChange: (v: string) => void;
}) {
  const refs = useRef<Array<HTMLInputElement | null>>([]),
    digits = Array.from({ length: 6 }, (_, i) => value[i] || "");
  return (
    <div className="otp-boxes">
      {digits.map((d, i) => (
        <input
          key={i}
          ref={(el) => {
            refs.current[i] = el;
          }}
          className="otp-box"
          inputMode="numeric"
          maxLength={1}
          value={d}
          required
          onChange={(e) => {
            const n = [...digits];
            n[i] = e.target.value.replace(/\D/g, "").slice(-1);
            onChange(n.join(""));
            if (n[i] && i < 5) refs.current[i + 1]?.focus();
          }}
          onKeyDown={(e) => {
            if (e.key === "Backspace" && !d && i > 0)
              refs.current[i - 1]?.focus();
          }}
        />
      ))}
    </div>
  );
}
function SOSSelect({
  options,
  value,
  onChange,
}: {
  options: any[];
  value: string;
  onChange: (value: string) => void;
}) {
  const [open, setOpen] = useState(false),
    selected = options.find((item) => String(item.id_sos) === value);
  return (
    <div className="sos-picker">
      <button
        type="button"
        className={open ? "sos-trigger open" : "sos-trigger"}
        aria-expanded={open}
        onClick={() => setOpen((v) => !v)}
      >
        <span>
          {selected ? (
            <>
              <b>{selected.name_sos}</b>
              <small>{selected.telepon}</small>
            </>
          ) : (
            <span className="picker-placeholder">Choose SOS contact</span>
          )}
        </span>
        <span className="picker-chevron">?</span>
      </button>
      {open && (
        <div className="sos-dropdown">
          <div className="sos-options">
            {options.length ? (
              options.map((item) => {
                const id = String(item.id_sos),
                  active = id === value;
                return (
                  <button
                    type="button"
                    className={active ? "sos-option selected" : "sos-option"}
                    key={id}
                    onClick={() => {
                      onChange(id);
                      setOpen(false);
                    }}
                  >
                    <span>
                      <b>{item.name_sos}</b>
                      <small>
                        {item.telepon}
                        {item.alamat_sos ? ` · ${item.alamat_sos}` : ""}
                      </small>
                    </span>
                    {active && <i>?</i>}
                  </button>
                );
              })
            ) : (
              <p className="facility-empty">No SOS contacts available</p>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
function FacilityMultiSelect({
  options,
  value,
  onChange,
}: {
  options: any[];
  value: string;
  onChange: (value: string) => void;
}) {
  const selected = value.split(",").filter(Boolean),
    [open, setOpen] = useState(false),
    [draft, setDraft] = useState<string[]>(selected);
  const toggle = (id: string) =>
    setDraft((current) =>
      current.includes(id)
        ? current.filter((item) => item !== id)
        : [...current, id],
    );
  const iconSrc = (icon: string) =>
    !icon
      ? ""
      : icon.startsWith("data:") || icon.startsWith("http")
        ? icon
        : `data:image/webp;base64,${icon}`;
  return (
    <div className="facility-picker">
      <button
        type="button"
        className={open ? "facility-trigger open" : "facility-trigger"}
        aria-expanded={open}
        onClick={() => {
          setDraft(selected);
          setOpen((v) => !v);
        }}
      >
        <span>
          {selected.length
            ? `${selected.length} facilities selected`
            : "Choose facilities"}
        </span>
        <span className="facility-chevron">?</span>
      </button>
      {selected.length > 0 && (
        <div className="selected-facilities">
          {selected.map((id) => {
            const item = options.find(
              (option) => String(option.id_facility) === id,
            );
            return item ? (
              <span className="facility-chip" key={id}>
                {item.icon ? (
                  <img src={iconSrc(item.icon)} alt="" />
                ) : (
                  <Icon name="facility" />
                )}
                <span>{item.namefacility}</span>
              </span>
            ) : null;
          })}
        </div>
      )}
      {open && (
        <div className="facility-dropdown">
          <div className="facility-options">
            {options.length ? (
              options.map((item) => {
                const id = String(item.id_facility),
                  checked = draft.includes(id);
                return (
                  <label
                    className={
                      checked ? "facility-option selected" : "facility-option"
                    }
                    key={id}
                  >
                    <input
                      type="checkbox"
                      checked={checked}
                      onChange={() => toggle(id)}
                    />
                    {item.icon ? (
                      <img src={iconSrc(item.icon)} alt="" />
                    ) : (
                      <span className="facility-placeholder">
                        <Icon name="facility" />
                      </span>
                    )}
                    <span>{item.namefacility}</span>
                    <i>{checked ? "?" : ""}</i>
                  </label>
                );
              })
            ) : (
              <p className="facility-empty">No facilities available</p>
            )}
          </div>
          <div className="facility-picker-actions">
            <button type="button" onClick={() => setOpen(false)}>
              Cancel
            </button>
            <button
              type="button"
              className="confirm"
              onClick={() => {
                onChange(draft.join(","));
                setOpen(false);
              }}
              disabled={!draft.length}
            >
              OK ({draft.length})
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
function SubcategoryMultiSelect({ options, value, onChange }: { options: any[]; value: string; onChange: (value: string) => void }) {
  const selected = value.split(",").filter(Boolean),
    [open, setOpen] = useState(false),
    [draft, setDraft] = useState<string[]>(selected);
  const toggle = (id: string) => setDraft((current) => current.includes(id) ? current.filter((item) => item !== id) : [...current, id]);
  return (
    <div className="subcategory-picker">
      <button type="button" className={open ? "subcategory-picker-trigger open" : "subcategory-picker-trigger"} aria-expanded={open} onClick={() => { setDraft(selected); setOpen((current) => !current); }}>
        <span>{selected.length ? `${selected.length} subcategories selected` : "Choose subcategories"}</span><span className="facility-chevron">?</span>
      </button>
      {selected.length > 0 && <div className="subcategory-picker-chips">{selected.map((id) => {
        const item = options.find((option) => String(option.id_subcategories) === id);
        return item ? <span key={id}><span>{item.namesubcategories}<small>{item.category?.name || "Uncategorized"}</small></span><button type="button" aria-label={`Remove ${item.namesubcategories}`} onClick={() => onChange(selected.filter((itemId) => itemId !== id).join(","))}><Icon name="close" /></button></span> : null;
      })}</div>}
      {open && <div className="subcategory-picker-dropdown">
        <div className="subcategory-picker-options">{options.length ? options.map((item) => {
          const id = String(item.id_subcategories), checked = draft.includes(id);
          return <label className={checked ? "selected" : ""} key={id}><input type="checkbox" checked={checked} onChange={() => toggle(id)} /><span><b>{item.namesubcategories}</b><small>{item.category?.name || "Uncategorized"}</small></span><i>{checked ? "?" : ""}</i></label>;
        }) : <p className="facility-empty">No subcategories available</p>}</div>
        <div className="facility-picker-actions"><button type="button" onClick={() => setOpen(false)}>Cancel</button><button type="button" className="confirm" disabled={!draft.length} onClick={() => { onChange(draft.join(",")); setOpen(false); }}>Apply ({draft.length})</button></div>
      </div>}
    </div>
  );
}
function Header({ user }: { user: any }) {
  return (
    <header>
      <span>Dashboard</span>
      <span className="admin-profile">
        <span>
          {user?.username || "Admin"}{" "}
          <b className="desktop-only">| Padmanabha Manager</b>
        </span>
        <i>
          <Icon name="user" />
        </i>
      </span>
    </header>
  );
}
function Decor() {
  return (
    <>
      <div className="auth-decoration auth-decoration-one" />
      <div className="auth-decoration auth-decoration-two" />
    </>
  );
}
function Logo() {
  return (
    <div className="auth-logo">
      <img src="/padmanabha-logo.png" alt="Padmanabha" />
    </div>
  );
}
function Heading({
  tag,
  title,
  text,
}: {
  tag: string;
  title: string;
  text: string;
}) {
  return (
    <div className="auth-heading">
      <span>{tag}</span>
      <h1>{title}</h1>
      <p>{text}</p>
    </div>
  );
}
function Submit({
  busy,
  busyText,
  text,
}: {
  busy: boolean;
  busyText: string;
  text: string;
}) {
  return (
    <button className="login-button" disabled={busy}>
      {busy ? (
        <>
          <span className="button-spinner" />
          {busyText}
        </>
      ) : (
        text
      )}
    </button>
  );
}
function Error({ text }: { text: string }) {
  return (
    <p className="auth-error">
      <span>!</span>
      {text}
    </p>
  );
}
function MailIcon() {
  return (
    <svg className="field-icon" viewBox="0 0 24 24">
      <path d="M4 6h16v12H4z" />
      <path d="m4 7 8 6 8-6" />
    </svg>
  );
}
function LockIcon() {
  return (
    <svg className="field-icon" viewBox="0 0 24 24">
      <rect x="5" y="10" width="14" height="11" rx="2" />
      <path d="M8 10V7a4 4 0 0 1 8 0v3" />
    </svg>
  );
}
function Eye({ off }: { off: boolean }) {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      {off ? (
        <>
          <path d="M3 3l18 18" />
          <path d="M10.6 6.2A10.4 10.4 0 0 1 12 6c6 0 9.5 6 9.5 6a16 16 0 0 1-2.1 2.8M6.2 6.2C3.8 8 2.5 12 2.5 12s3.5 6 9.5 6c1.5 0 2.8-.4 4-.9M9.9 9.9a3 3 0 0 0 4.2 4.2" />
        </>
      ) : (
        <>
          <path d="M2.5 12s3.5-6 9.5-6 9.5 6 9.5 6-3.5 6-9.5 6-9.5-6-9.5-6Z" />
          <circle cx="12" cy="12" r="2.5" />
        </>
      )}
    </svg>
  );
}

function Icon({ name }: { name: string }) {
  const paths: Record<string, ReactNode> = {
    overview: (
      <>
        <rect x="3" y="3" width="7" height="7" rx="1" />
        <rect x="14" y="3" width="7" height="7" rx="1" />
        <rect x="3" y="14" width="7" height="7" rx="1" />
        <rect x="14" y="14" width="7" height="7" rx="1" />
      </>
    ),
    villa: (
      <>
        <path d="M12 21s6-5.2 6-11a6 6 0 1 0-12 0c0 5.8 6 11 6 11Z" />
        <circle cx="12" cy="10" r="2" />
      </>
    ),
    activity: (
      <>
        <path d="M6 18 9.5 7l2.5 6 2-4 4 9" />
        <path d="M4 18h16M15 6h.01" />
      </>
    ),
    restaurant: (
      <>
        <path d="M6 3v8M3 3v5a3 3 0 0 0 6 0V3M6 11v10" />
        <path d="M15 3v18M15 3c4 1 5 4 5 8h-5" />
      </>
    ),
    event: (
      <>
        <rect x="3" y="5" width="18" height="16" rx="2" />
        <path d="M16 3v4M8 3v4M3 10h18" />
      </>
    ),
    category: <path d="M4 4h6v6H4zM14 4h6v6h-6zM4 14h6v6H4zM14 14h6v6h-6z" />,
    subcategory: (
      <>
        <path d="M5 5h5v5H5zM14 5h5v5h-5zM14 14h5v5h-5z" />
        <path d="M7.5 10v6.5H14M10 7.5h4" />
      </>
    ),
    user: (
      <>
        <circle cx="12" cy="8" r="4" />
        <path d="M4 21a8 8 0 0 1 16 0" />
      </>
    ),
    sos: (
      <>
        <circle cx="12" cy="12" r="9" />
        <path d="M12 7v6M12 17h.01" />
      </>
    ),
    facility: <path d="M4 21V5l8-3 8 3v16M10 21v-4h4v4" />,
    more: (
      <>
        <circle cx="5" cy="12" r="1" fill="currentColor" stroke="none" />
        <circle cx="12" cy="12" r="1" fill="currentColor" stroke="none" />
        <circle cx="19" cy="12" r="1" fill="currentColor" stroke="none" />
      </>
    ),
    plus: <path d="M12 5v14M5 12h14" />,
    search: (
      <>
        <circle cx="11" cy="11" r="7" />
        <path d="m20 20-4-4" />
      </>
    ),
    edit: (
      <>
        <path d="m4 20 4.2-1 10.7-10.7a2 2 0 0 0-2.8-2.8L5.4 16.2 4 20Z" />
      </>
    ),
    delete: (
      <>
        <path d="M4 7h16M9 7V4h6v3M7 7l1 14h8l1-14" />
      </>
    ),
    eye: (
      <>
        <path d="M2.5 12s3.5-6 9.5-6 9.5 6 9.5 6-3.5 6-9.5 6-9.5-6-9.5-6Z" />
        <circle cx="12" cy="12" r="2.5" />
      </>
    ),
    logout: (
      <>
        <path d="M10 17l5-5-5-5M15 12H3" />
        <path d="M14 4h5a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2h-5" />
      </>
    ),
    menu: <path d="M4 7h16M4 12h16M4 17h16" />,
    chevrondown: <path d="m7 10 5 5 5-5" />,
    close: <path d="m6 6 12 12M18 6 6 18" />,
    check: <path d="m5 12 4 4L19 6" />,
    arrowleft: (
      <>
        <path d="m10 6-6 6 6 6M4 12h16" />
      </>
    ),
    image: (
      <>
        <rect x="3" y="4" width="18" height="16" rx="2" />
        <path d="m3 17 5-5 4 4 3-3 6 6" />
      </>
    ),
  };
  return (
    <svg className="ui-icon" viewBox="0 0 24 24" aria-hidden="true">
      {paths[name] || paths.category}
    </svg>
  );
}

function useRows(path: string) {
  const [rows, setRows] = useState<any[]>([]),
    [loading, setLoading] = useState(true),
    [error, setError] = useState("");
  const load = async () => {
    setLoading(true);
    setError("");
    try {
      setRows(unwrap(await api(path)));
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };
  useEffect(() => {
    load();
  }, [path]);
  return { rows, setRows, loading, error, load };
}

type ResourceConfig = {
  title: string;
  path: string;
  columns: { label: string; value: (row: any) => ReactNode }[];
  add?: boolean;
};
function ResourcePage({ config }: { config: ResourceConfig }) {
  const data = useRows(config.path);
  const [query, setQuery] = useState(""),
    [editing, setEditing] = useState<any>(null),
    [show, setShow] = useState(false),
    [detail, setDetail] = useState<any>(null);
  const usesPopupEditor = ["/admin/villa", "/admin/activity", "/admin/restaurant", "/admin/event", "/admin/facility", "/admin/sos"].includes(config.path);
  const filtered = useMemo(
    () =>
      data.rows.filter((r) =>
        JSON.stringify(r).toLowerCase().includes(query.toLowerCase()),
      ),
    [data.rows, query],
  );
  async function remove(row: any) {
    if (!confirm("Are you sure you want to delete this item?")) return;
    try {
      await api(`${config.path}/${idOf(row)}`, { method: "DELETE" });
      data.setRows((v) => v.filter((x) => idOf(x) !== idOf(row)));
    } catch (e: any) {
      alert(e.message);
    }
  }
  return (
    <section>
      <div className="title-row">
        <div className="page-title">
          <span>Content management</span>
          <h2>{config.title}</h2>
        </div>
        <div className="page-actions">
          <div className="toolbar">
            <input
              className="search"
              placeholder="Search data..."
              value={query}
              onChange={(e) => setQuery(e.target.value)}
            />
          </div>
          {config.add && (
            <button
              className="primary"
              onClick={() => {
                setEditing(null);
                setShow(true);
              }}
            >
              <Icon name="plus" />
              <span>
                Add{" "}
                {config.title.replace(/^Manage (Your )?/, "").replace(/s$/, "")}
              </span>
            </button>
          )}
        </div>
      </div>
      {show && (
        usesPopupEditor ? (
          <div className={`modal resource-editor-overlay${config.path === "/admin/facility" ? " facility-editor-overlay" : config.path === "/admin/sos" ? " sos-editor-overlay" : config.path === "/admin/restaurant" ? " restaurant-editor-overlay" : config.path === "/admin/activity" ? " activity-editor-overlay" : config.path === "/admin/villa" ? " villa-editor-overlay" : ""}`} role="presentation" onMouseDown={(e) => { if (e.target === e.currentTarget) setShow(false); }}>
            <Editor config={config} value={editing} onClose={() => setShow(false)} onSaved={() => { setShow(false); data.load(); }} />
          </div>
        ) : (
          <Editor config={config} value={editing} onClose={() => setShow(false)} onSaved={() => { setShow(false); data.load(); }} />
        )
      )}
      {detail && <Detail value={detail} onClose={() => setDetail(null)} />}{" "}
      {data.loading ? (
        <div className="spinner" />
      ) : data.error ? (
        <div className="error-box">
          <p>{data.error}</p>
          <button onClick={data.load}>Try again</button>
        </div>
      ) : (
        <div className="card table-wrap">
          <table>
            <thead>
              <tr>
                <th className="number-column">No.</th>
                {config.columns.map((c) => (
                  <th key={c.label}>{c.label}</th>
                ))}
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {filtered.map((r, i) => (
                <tr key={idOf(r) ?? i}>
                  <td className="number-column">{i + 1}</td>
                  {config.columns.map((c) => (
                    <td key={c.label}>{c.value(r) ?? "-"}</td>
                  ))}
                  <td className="actions">
                    {["villa", "activity", "restaurant", "event", "package"].includes(
                      config.path.split("/").pop() || "",
                    ) && (
                      <button
                        className="action-button detail-action"
                        onClick={() => setDetail(r)}
                        title="View details"
                      >
                        <Icon name="eye" />
                        <span>Detail</span>
                      </button>
                    )}
                    {config.add && (
                      <button
                        className="action-button edit-action"
                        onClick={() => {
                          setEditing(r);
                          setShow(true);
                        }}
                        title="Edit"
                      >
                        <Icon name="edit" />
                        <span>Edit</span>
                      </button>
                    )}
                    <button
                      className="action-button danger"
                      onClick={() => remove(r)}
                      title="Delete"
                    >
                      <Icon name="delete" />
                      <span>Delete</span>
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {!filtered.length && <p className="empty">No data found</p>}
        </div>
      )}
    </section>
  );
}
function Editor({
  config,
  value,
  onClose,
  onSaved,
}: {
  config: ResourceConfig;
  value: any;
  onClose: () => void;
  onSaved: () => void;
}) {
  const type = config.path.split("/").pop()!;
  const [facilityOptions, setFacilityOptions] = useState<any[]>([]),
    [subcategoryOptions, setSubcategoryOptions] = useState<any[]>([]);
  useEffect(() => {
    if (["villa", "activity", "restaurant"].includes(type)) {
      Promise.all([
        api("/admin/facility"),
        api("/admin/subcategory"),
      ])
        .then(([facilities, subcategories]) => {
          setFacilityOptions(unwrap(facilities));
          setSubcategoryOptions(unwrap(subcategories));
        })
        .catch(() => {});
    }
  }, [type]);
  const fields: Record<string, [string, string, string][]> = {
    villa: [
      ["namevilla", "Villa name", "text"],
      ["subcategoryId", "Subcategories", "subcategory-select"],
      ["facilityId", "Facilities", "facility-select"],
      ["operational", "Operational hours", "operational-hours"],
      ["description", "Description", "textarea"],
      ["image", "Villa images", "files"],
    ],
    activity: [
      ["nameactivity", "Activity name", "text"],
      ["subcategoryId", "Subcategories", "subcategory-select"],
      ["facilityId", "Facilities", "facility-select"],
      ["operational", "Operational hours", "operational-hours"],
      ["description", "Description", "textarea"],
      ["image", "Activity images", "files"],
    ],
    restaurant: [
      ["namerestaurant", "Restaurant name", "text"],
      ["subcategoryId", "Subcategories", "subcategory-select"],
      ["facilityId", "Facilities", "facility-select"],
      ["operational", "Operational hours", "operational-hours"],
      ["description", "Description", "textarea"],
      ["image", "Restaurant images", "files"],
    ],
    event: [
      ["nameevent", "Event name", "text"],
      ["price", "Price", "number"],
      ["start_date", "Start date", "date"],
      ["end_date", "End date", "date"],
      ["start_time", "Start time", "time"],
      ["end_time", "End time", "time"],
      ["description", "Description", "textarea"],
      ["image", "Event images", "files"],
    ],
    sos: [
      ["name_sos", "Name", "text"],
      ["alamat_sos", "Location", "text"],
      ["telepon", "Number", "text"],
    ],
    facility: [
      ["namefacility", "Name", "text"],
      ["icon", "Icon", "file"],
    ],
  };
  const list = fields[type] || [];
  const entityLabel = type === "villa" ? "Villa" : type === "activity" ? "Activity" : type === "restaurant" ? "Restaurant" : type === "event" ? "Event" : type === "facility" ? "Facility" : type === "sos" ? "SOS" : "Data";
  const requiredKeys: Record<string, string[]> = {
    villa: ["namevilla", "subcategoryId", "facilityId"],
    activity: ["nameactivity", "subcategoryId", "facilityId"],
    restaurant: ["namerestaurant", "subcategoryId", "facilityId"],
    event: ["nameevent", "start_date", "end_date", "start_time", "end_time"],
    sos: ["name_sos", "alamat_sos", "telepon"],
  };
  const isRequiredField = (key: string, fieldType: string) =>
    (requiredKeys[type] || []).includes(key) || ((fieldType === "files" || fieldType === "file") && !value);
  const [form, setForm] = useState<Record<string, any>>(() => {
      const initial = Object.fromEntries(list.map(([k]) => [k, value?.[k] ?? ""]));
      return initial;
    }),
    [busy, setBusy] = useState(false),
    [error, setError] = useState("");
  async function save(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    const missing = list.find(([key, , fieldType]) => {
      if (!isRequiredField(key, fieldType)) return false;
      const fieldValue = form[key];
      return Array.isArray(fieldValue) ? fieldValue.length === 0 : !String(fieldValue || "").trim();
    });
    if (missing) {
      setError(`${missing[1]} is required.`);
      return;
    }
    setBusy(true);
    try {
      const hasFile = list.some(([, , t]) => t === "file" || t === "files");
      let body: BodyInit;
      if (hasFile) {
        const fd = new FormData();
        Object.entries(form).forEach(([k, v]) => {
          if (v instanceof File) fd.append(k, v);
          else if (Array.isArray(v)) v.forEach((file) => fd.append(k, file));
          else if (v !== "") fd.append(k, String(v));
        });
        body = fd;
      } else body = JSON.stringify(form);
      await api(value ? `${config.path}/${idOf(value)}` : config.path, {
        method: value ? "PUT" : "POST",
        body,
      });
      onSaved();
    } catch (e: any) {
      setError(e.message);
    } finally {
      setBusy(false);
    }
  }
  return (
    <div className="card editor resource-editor-card">
      <div className="editor-head">
        <div><span>{value ? "Update information" : `New ${entityLabel.toLowerCase()}`}</span><h3>{value ? "Edit" : "Add"} {entityLabel}</h3></div>
        <button type="button" aria-label="Close editor" onClick={onClose}>
          <Icon name="close" />
        </button>
      </div>
      <form onSubmit={save} className="form-grid">
        {list.map(([key, label, type]) => (
          <label key={key}>
            <span className="field-label"><span>{label}{isRequiredField(key, type) && <b className="required-mark" aria-label="required">*</b>}</span></span>
            {type === "textarea" ? (
              <textarea
                value={form[key]}
                placeholder={`Enter ${label.toLowerCase()}...`}
                onChange={(e) => setForm({ ...form, [key]: e.target.value })}
              />
            ) : type === "facility-select" ? (
              <FacilityMultiSelect
                options={facilityOptions}
                value={String(form[key] || "")}
                onChange={(selected) => setForm({ ...form, [key]: selected })}
              />
            ) : type === "subcategory-select" ? (
              <SubcategoryMultiSelect options={subcategoryOptions} value={String(form[key] || "")} onChange={(selected) => setForm({ ...form, [key]: selected })} />
            ) : type === "operational-hours" ? (
              <div className="operational-hours-input">
                <span className="operational-time-field">
                  <input
                    type="time"
                    aria-label="Opening time"
                    className={String(form[key] || "").split(" - ")[0] ? "has-value" : ""}
                    value={String(form[key] || "").split(" - ")[0] || ""}
                    onClick={(e) => e.currentTarget.showPicker?.()}
                    onChange={(e) => {
                      const closingTime = String(form[key] || "").split(" - ")[1] || "";
                      setForm({ ...form, [key]: `${e.target.value} - ${closingTime}` });
                    }}
                  />
                  {!String(form[key] || "").split(" - ")[0] && <small>Opening time</small>}
                </span>
                <i>to</i>
                <span className="operational-time-field">
                  <input
                    type="time"
                    aria-label="Closing time"
                    className={String(form[key] || "").split(" - ")[1] ? "has-value" : ""}
                    value={String(form[key] || "").split(" - ")[1] || ""}
                    onClick={(e) => e.currentTarget.showPicker?.()}
                    onChange={(e) => {
                      const openingTime = String(form[key] || "").split(" - ")[0] || "";
                      setForm({ ...form, [key]: `${openingTime} - ${e.target.value}` });
                    }}
                  />
                  {!String(form[key] || "").split(" - ")[1] && <small>Closing time</small>}
                </span>
              </div>
            ) : type === "files" ? (
              <input
                type="file"
                accept="image/*"
                multiple
                required={!value}
                onChange={(e) =>
                  setForm({ ...form, [key]: Array.from(e.target.files || []) })
                }
              />
            ) : key === "price" ? (
              <div className="currency-input">
                <span>Rp</span>
                <input
                  type="number"
                  min="0"
                  step="any"
                  placeholder="0"
                  required={isRequiredField(key, type)}
                  defaultValue={form[key]}
                  onChange={(e) => setForm({ ...form, [key]: e.target.value })}
                />
              </div>
            ) : (
              <input
                type={type}
                placeholder={type === "url" ? "https://maps.google.com/..." : `Enter ${label.toLowerCase()}...`}
                step={type === "number" ? "any" : undefined}
                required={isRequiredField(key, type)}
                defaultValue={type === "file" ? undefined : form[key]}
                onClick={(e) => {
                  if (type === "date" || type === "time") e.currentTarget.showPicker?.();
                }}
                onChange={(e) =>
                  setForm({
                    ...form,
                    [key]:
                      type === "file" ? e.target.files?.[0] : e.target.value,
                  })
                }
              />
            )}
          </label>
        ))}
        {error && <p className="error">{error}</p>}
        <button className="primary save" disabled={busy}>
          {busy ? "Saving..." : "Save"}
        </button>
      </form>
    </div>
  );
}
function Detail({ value, onClose }: { value: any; onClose: () => void }) {
  return (
    <div className="modal">
      <div className="card detail">
        <button className="modal-close" onClick={onClose}>
          <Icon name="close" />
        </button>
        <h2>
          {value.namevilla ||
            value.nameactivity ||
            value.namerestaurant ||
            value.nameevent ||
            value.villa?.namevilla ||
            "Detail"}
        </h2>
        {Object.entries(value)
          .filter(([k, v]) => !k.startsWith("id_") && !["do", "dont", "maps", "safety", "sosId"].includes(k) && typeof v !== "object")
          .map(([k, v]) => (
            <div className="detail-row" key={k}>
              <b>{k.replaceAll("_", " ")}</b>
              <span>{String(v ?? "-")}</span>
            </div>
          ))}
      </div>
    </div>
  );
}

function ModulePlaceholder({ title, description, icon }: { title: string; description: string; icon: string }) {
  return (
    <section className="module-placeholder-page">
      <div className="title-row categories-heading">
        <div className="page-title"><span>Content management</span><h2>{title}</h2><p>{description}</p></div>
      </div>
      <div className="module-placeholder card">
        <span><Icon name={icon} /></span>
        <h3>{title} module</h3>
        <p>This menu is ready. Data management will become available when its backend API is connected.</p>
      </div>
    </section>
  );
}

function OverviewPage() {
  const [counts, setCounts] = useState([0, 0, 0, 0]),
    [loading, setLoading] = useState(true);
  useEffect(() => {
    Promise.all(
      [
        "/admin/villa",
        "/admin/event",
        "/admin/category",
        "/admin/users",
      ].map((p) => api(p).then((x) => unwrap(x).length)),
    )
      .then(setCounts)
      .finally(() => setLoading(false));
  }, []);
  return (
    <section className="overview">
      {loading ? (
        <div className="spinner" />
      ) : (
        <div className="stats">
          {["Villas", "Events", "Categories", "Users active"].map(
            (x, i) => (
              <div className="card stat" key={x}>
                <span>{x}</span>
                <strong>{counts[i]}</strong>
              </div>
            ),
          )}
        </div>
      )}
    </section>
  );
}

function VillasPage() {
  return (
    <ResourcePage
      config={{
        title: "Manage Villas",
        path: "/admin/villa",
        add: true,
        columns: [
          { label: "Name", value: (r) => r.namevilla },
          {
            label: "Category",
            value: (r) => r.subcategory?.[0]?.category?.name || "-",
          },
          {
            label: "Sub Category",
            value: (r) =>
              r.subcategory
                ?.map((x: any) => x.namesubcategories || x.name)
                .join(", ") || "-",
          },
        ],
      }}
    />
  );
}

function ActivitiesPage() {
  return <ResourcePage config={{ title: "Manage Activities", path: "/admin/activity", add: true, columns: [
    { label: "Name", value: (r) => r.nameactivity },
    { label: "Category", value: (r) => r.subcategory?.[0]?.category?.name || "-" },
    { label: "Sub Category", value: (r) => r.subcategory?.map((x: any) => x.namesubcategories || x.name).join(", ") || "-" },
  ] }} />;
}

function RestaurantsPage() {
  return <ResourcePage config={{ title: "Manage Restaurants", path: "/admin/restaurant", add: true, columns: [
    { label: "Name", value: (r) => r.namerestaurant },
    { label: "Category", value: (r) => r.subcategory?.[0]?.category?.name || "-" },
    { label: "Sub Category", value: (r) => r.subcategory?.map((x: any) => x.namesubcategories || x.name).join(", ") || "-" },
  ] }} />;
}

function EventsPage() {
  return (
    <ResourcePage
      config={{
        title: "Manage Events",
        path: "/admin/event",
        add: true,
        columns: [
          { label: "Name", value: (r) => r.nameevent },
          {
            label: "Date",
            value: (r) =>
              `${r.start_date || "-"}${r.end_date ? ` - ${r.end_date}` : ""}`,
          },
          {
            label: "Price",
            value: (r) =>
              r.price == null
                ? "-"
                : `Rp ${Number(r.price).toLocaleString("id-ID")}`,
          },
        ],
      }}
    />
  );
}

function CategoriesPage() {
  const cats = useRows("/admin/category");
  const [name, setName] = useState(""),
    [query, setQuery] = useState(""),
    [showCreate, setShowCreate] = useState(false),
    [editingCategory, setEditingCategory] = useState<any>(null),
    [editCategoryName, setEditCategoryName] = useState("");
  async function add(path: string, body: any) {
    try {
      await api(path, { method: "POST", body: JSON.stringify(body) });
      await cats.load();
    } catch (e: any) {
      alert(e.message);
    }
  }
  async function del(path: string, id: any) {
    if (confirm("Delete this item? This action cannot be undone.")) {
      await api(`${path}/${id}`, { method: "DELETE" });
      await cats.load();
    }
  }
  async function updateCategory(e: React.FormEvent) {
    e.preventDefault();
    if (!editingCategory || !editCategoryName.trim()) return;
    try {
      await api(`/admin/category/${editingCategory.id_categories}`, {
        method: "PUT",
        body: JSON.stringify({ name: editCategoryName.trim() }),
      });
      setEditingCategory(null);
      setEditCategoryName("");
      await cats.load();
    } catch (e: any) {
      alert(e.message);
    }
  }
  const shown = cats.rows.filter((x) =>
    x.name?.toLowerCase().includes(query.toLowerCase()),
  );
  return (
    <section className="categories-page">
      <div className="title-row categories-heading">
        <div className="page-title">
          <span>Content management</span>
          <h2>Categories</h2>
          <p>Create and manage your content categories.</p>
        </div>
        <div className="category-summary" aria-label="Category summary">
          <span><b>{cats.rows.length}</b> Categories</span>
        </div>
      </div>
      <div className="category-controls">
        <div className="category-search">
          <Icon name="search" />
          <input placeholder="Search categories..." aria-label="Search categories" value={query} onChange={(e) => setQuery(e.target.value)} />
        </div>
        <button className="primary category-add-trigger" onClick={() => setShowCreate(true)}>
          <Icon name="plus" />
          <span>Add Category</span>
        </button>
      </div>
      {showCreate && (
        <div className="modal category-modal" role="presentation" onMouseDown={(e) => {
          if (e.target === e.currentTarget) setShowCreate(false);
        }}>
          <form className="category-modal-card" role="dialog" aria-modal="true" aria-labelledby="new-category-title" onSubmit={(e) => {
            e.preventDefault();
            if (name.trim()) add("/admin/category", { name: name.trim() }).then(() => {
              setName("");
              setShowCreate(false);
            });
          }}>
            <div className="category-modal-head">
              <div><span>New category</span><h3 id="new-category-title">Add Category</h3></div>
              <button type="button" aria-label="Close popup" onClick={() => setShowCreate(false)}><Icon name="close" /></button>
            </div>
            <div className="category-modal-body">
              <label htmlFor="category-name">Category name</label>
              <input id="category-name" autoFocus placeholder="e.g. Adventure" value={name} onChange={(e) => setName(e.target.value)} />
              <p>After creating it, you can assign subcategories from the Sub Categories page.</p>
            </div>
            <div className="category-modal-actions">
              <button type="button" className="category-cancel" onClick={() => setShowCreate(false)}>Cancel</button>
              <button className="primary" type="submit" disabled={!name.trim()}><Icon name="plus" /><span>Add Category</span></button>
            </div>
          </form>
        </div>
      )}
      {editingCategory && (
        <div className="modal category-modal" role="presentation" onMouseDown={(e) => { if (e.target === e.currentTarget) setEditingCategory(null); }}>
          <form className="category-modal-card" role="dialog" aria-modal="true" aria-labelledby="edit-category-title" onSubmit={updateCategory}>
            <div className="category-modal-head"><div><span>Update category</span><h3 id="edit-category-title">Edit Category</h3></div><button type="button" aria-label="Close popup" onClick={() => setEditingCategory(null)}><Icon name="close" /></button></div>
            <div className="category-modal-body"><label htmlFor="edit-category-name">Category name</label><input id="edit-category-name" autoFocus value={editCategoryName} onChange={(e) => setEditCategoryName(e.target.value)} required /><p>Update the category name, then save your changes.</p></div>
            <div className="category-modal-actions"><button type="button" className="category-cancel" onClick={() => setEditingCategory(null)}>Cancel</button><button className="primary" type="submit" disabled={!editCategoryName.trim()}><Icon name="check" /><span>Save Changes</span></button></div>
          </form>
        </div>
      )}
      {cats.loading ? <div className="spinner" /> : cats.error ? (
        <div className="error-box card"><p>{cats.error}</p><button onClick={cats.load}>Try again</button></div>
      ) : shown.length ? (
        <div className="category-list">
          {shown.map((cat) => <CategoryCard key={cat.id_categories} cat={cat} del={del} onEdit={() => { setEditingCategory(cat); setEditCategoryName(cat.name || ""); }} />)}
        </div>
      ) : (
        <div className="category-empty card"><span><Icon name="category" /></span><h3>{query ? "No matching categories" : "No categories yet"}</h3><p>{query ? "Try another search term." : "Add your first category using the form above."}</p></div>
      )}
    </section>
  );
}
function CategoryCard({ cat, del, onEdit }: any) {
  return (
    <article className="category-card card">
      <div className="category-card-head">
        <div className="category-card-title"><span><Icon name="category" /></span><div><h3>{cat.name}</h3><p>Category</p></div></div>
        <div className="category-card-actions"><button className="edit-action" type="button" onClick={onEdit}><Icon name="edit" /><span>Edit</span></button><button className="danger" type="button" onClick={() => del("/admin/category", cat.id_categories)}><Icon name="delete" /><span>Remove</span></button></div>
      </div>
    </article>
  );
}

function SubCategoriesPage() {
  const data = useRows("/admin/subcategory"),
    categories = useRows("/admin/category");
  const [query, setQuery] = useState(""),
    [showCreate, setShowCreate] = useState(false),
    [name, setName] = useState(""),
    [categoryId, setCategoryId] = useState(""),
    [selectedNames, setSelectedNames] = useState<string[]>([]),
    [nameError, setNameError] = useState(""),
    [detailItem, setDetailItem] = useState<any>(null),
    [editingItem, setEditingItem] = useState<any>(null),
    [editName, setEditName] = useState(""),
    [editCategoryId, setEditCategoryId] = useState("");
  const subcategoryDropdownRef = useRef<HTMLDetailsElement>(null);
  const normalizedQuery = query.trim().toLowerCase();
  const capitalizeName = (value: string) => {
    const clean = value.trim();
    return clean ? clean.charAt(0).toUpperCase() + clean.slice(1) : "";
  };
  const assignedNames = new Set(
    data.rows
      .filter((item) => String(item.category?.id_categories ?? item.categoriesId) === categoryId)
      .map((item) => String(item.namesubcategories || "").toLowerCase()),
  );
  const existingNameOptions = Array.from(
    new Set(data.rows.map((item) => capitalizeName(item.namesubcategories || "")).filter(Boolean)),
  ).filter((item) => !assignedNames.has(item.toLowerCase()));
  const groupedCategories = categories.rows.map((category) => {
    const allItems = data.rows.filter((item) =>
      String(item.category?.id_categories ?? item.categoriesId) === String(category.id_categories),
    );
    const categoryMatches = category.name?.toLowerCase().includes(normalizedQuery);
    const items = !normalizedQuery || categoryMatches
      ? allItems
      : allItems.filter((item) => item.namesubcategories?.toLowerCase().includes(normalizedQuery));
    return { category, items, categoryMatches };
  }).filter((group) => !normalizedQuery || group.categoryMatches || group.items.length);
  async function create(e: React.FormEvent) {
    e.preventDefault();
    const manualName = capitalizeName(name);
    const canonicalManualName = existingNameOptions.find((item) => item.toLowerCase() === manualName.toLowerCase()) || manualName;
    const names = Array.from(new Set([...selectedNames, ...(canonicalManualName ? [canonicalManualName] : [])]))
      .filter((item) => !assignedNames.has(item.toLowerCase()));
    if (!names.length || !categoryId) return;
    try {
      await Promise.all(names.map((subcategoryName) => api("/admin/subcategory", {
        method: "POST",
        body: JSON.stringify({
          namesubcategories: subcategoryName,
          categoriesId: Number(categoryId),
        }),
      })));
      setName("");
      setNameError("");
      setSelectedNames([]);
      setCategoryId("");
      setShowCreate(false);
      await data.load();
    } catch (e: any) {
      alert(e.message);
    }
  }
  async function remove(item: any) {
    if (!confirm(`Delete subcategory "${item.namesubcategories}"?`)) return;
    try {
      await api(`/admin/subcategory/${item.id_subcategories}`, { method: "DELETE" });
      await data.load();
    } catch (e: any) {
      alert(e.message);
    }
  }
  function addManualName() {
    const next = capitalizeName(name);
    if (!next) return;
    if (assignedNames.has(next.toLowerCase())) {
      setNameError(`"${next}" already exists in this category.`);
      return;
    }
    const canonicalName = existingNameOptions.find((item) => item.toLowerCase() === next.toLowerCase()) || next;
    setSelectedNames((current) => current.some((item) => item.toLowerCase() === canonicalName.toLowerCase()) ? current : [...current, canonicalName]);
    setName("");
    setNameError("");
  }
  async function updateSubcategory(e: React.FormEvent) {
    e.preventDefault();
    if (!editingItem || !editName.trim() || !editCategoryId) return;
    try {
      await api(`/admin/subcategory/${editingItem.id_subcategories}`, {
        method: "PUT",
        body: JSON.stringify({ namesubcategories: editName.trim(), categoriesId: Number(editCategoryId) }),
      });
      setEditingItem(null);
      await data.load();
    } catch (e: any) {
      alert(e.message);
    }
  }
  return (
    <section className="categories-page">
      <div className="title-row categories-heading">
        <div className="page-title">
          <span>Content management</span>
          <h2>Sub Categories</h2>
          <p>Manage subcategories and assign each one to a parent category.</p>
        </div>
        <div className="category-summary"><span><b>{data.rows.length}</b> Subcategories</span><i /><span><b>{categories.rows.length}</b> Categories</span></div>
      </div>
      <div className="category-controls">
        <div className="category-search"><Icon name="search" /><input placeholder="Search subcategories..." aria-label="Search subcategories" value={query} onChange={(e) => setQuery(e.target.value)} /></div>
        <button className="primary category-add-trigger" onClick={() => setShowCreate(true)}><Icon name="plus" /><span>Add Subcategory</span></button>
      </div>
      {showCreate && (
        <div className="modal category-modal" role="presentation" onMouseDown={(e) => { if (e.target === e.currentTarget) setShowCreate(false); }}>
          <form className="category-modal-card" role="dialog" aria-modal="true" aria-labelledby="new-subcategory-title" onSubmit={create}>
            <div className="category-modal-head"><div><span>New subcategory</span><h3 id="new-subcategory-title">Add Subcategory</h3></div><button type="button" aria-label="Close popup" onClick={() => setShowCreate(false)}><Icon name="close" /></button></div>
            <div className="category-modal-body subcategory-modal-fields">
              <label htmlFor="subcategory-parent">Parent category</label>
              <select id="subcategory-parent" value={categoryId} onChange={(e) => { setCategoryId(e.target.value); setSelectedNames([]); setName(""); setNameError(""); }} required>
                <option value="">Select a category</option>
                {categories.rows.map((category) => <option key={category.id_categories} value={category.id_categories}>{category.name}</option>)}
              </select>
              <label>Choose existing subcategories</label>
              <details className="subcategory-dropdown" ref={subcategoryDropdownRef}>
                <summary>{selectedNames.length ? `${selectedNames.length} selected` : "Select one or more subcategories"}<span>?</span></summary>
                <div className="subcategory-dropdown-menu">
                  {!categoryId ? <p>Select a parent category first.</p> : existingNameOptions.length ? existingNameOptions.map((item) => {
                    const checked = selectedNames.includes(item);
                    return <label className={checked ? "selected" : ""} key={item}><input type="checkbox" checked={checked} onChange={() => setSelectedNames((current) => current.includes(item) ? current.filter((value) => value !== item) : [...current, item])} /><span>{item}</span>{checked && <b>?</b>}</label>;
                  }) : <p>No other subcategories available.</p>}
                  {categoryId && <button className="subcategory-dropdown-done" type="button" onClick={() => subcategoryDropdownRef.current?.removeAttribute("open")}>Done</button>}
                </div>
              </details>
              <label htmlFor="subcategory-name">Or add a new subcategory</label>
              <div className={nameError ? "manual-subcategory-entry invalid" : "manual-subcategory-entry"}><input id="subcategory-name" placeholder="e.g. luxury" value={name} onChange={(e) => { const value = e.target.value; setName(value ? value.charAt(0).toUpperCase() + value.slice(1) : ""); setNameError(""); }} onKeyDown={(e) => { if (e.key === "Enter" && name.trim()) { e.preventDefault(); addManualName(); } }} /><button type="button" disabled={!name.trim() || !categoryId} onClick={addManualName}><Icon name="plus" /> Add</button></div>
              {nameError && <span className="subcategory-name-error">{nameError}</span>}
              {selectedNames.length > 0 && <div className="selected-subcategory-names">{selectedNames.map((item) => <span key={item}>{item}<button type="button" aria-label={`Remove ${item}`} onClick={() => setSelectedNames((current) => current.filter((value) => value !== item))}><Icon name="close" /></button></span>)}</div>}
              <p>Select existing names or type a new one. The first letter is capitalized automatically.</p>
            </div>
            <div className="category-modal-actions"><button type="button" className="category-cancel" onClick={() => setShowCreate(false)}>Cancel</button><button className="primary" type="submit" disabled={(!name.trim() && !selectedNames.length) || !categoryId}><Icon name="plus" /><span>Add {selectedNames.length > 1 ? `${selectedNames.length} Subcategories` : "Subcategory"}</span></button></div>
          </form>
        </div>
      )}
      {detailItem && (
        <div className="modal category-modal" role="presentation" onMouseDown={(e) => { if (e.target === e.currentTarget) setDetailItem(null); }}>
          <div className="category-modal-card" role="dialog" aria-modal="true" aria-labelledby="subcategory-detail-title">
            <div className="category-modal-head"><div><span>Subcategory information</span><h3 id="subcategory-detail-title">Subcategory Detail</h3></div><button type="button" aria-label="Close popup" onClick={() => setDetailItem(null)}><Icon name="close" /></button></div>
            <div className="subcategory-detail-body">
              <div><span>Name</span><b>{detailItem.namesubcategories}</b></div>
              <div><span>Parent category</span><b>{detailItem.category?.name || categories.rows.find((category) => String(category.id_categories) === String(detailItem.categoriesId))?.name || "—"}</b></div>
            </div>
            <div className="category-modal-actions"><button type="button" className="category-cancel" onClick={() => setDetailItem(null)}>Close</button></div>
          </div>
        </div>
      )}
      {editingItem && (
        <div className="modal category-modal" role="presentation" onMouseDown={(e) => { if (e.target === e.currentTarget) setEditingItem(null); }}>
          <form className="category-modal-card" role="dialog" aria-modal="true" aria-labelledby="edit-subcategory-title" onSubmit={updateSubcategory}>
            <div className="category-modal-head"><div><span>Update subcategory</span><h3 id="edit-subcategory-title">Edit Subcategory</h3></div><button type="button" aria-label="Close popup" onClick={() => setEditingItem(null)}><Icon name="close" /></button></div>
            <div className="category-modal-body subcategory-modal-fields">
              <label htmlFor="edit-subcategory-parent">Parent category</label>
              <select id="edit-subcategory-parent" value={editCategoryId} onChange={(e) => setEditCategoryId(e.target.value)} required>{categories.rows.map((category) => <option key={category.id_categories} value={category.id_categories}>{category.name}</option>)}</select>
              <label htmlFor="edit-subcategory-name">Subcategory name</label>
              <input id="edit-subcategory-name" autoFocus value={editName} onChange={(e) => setEditName(e.target.value)} required />
            </div>
            <div className="category-modal-actions"><button type="button" className="category-cancel" onClick={() => setEditingItem(null)}>Cancel</button><button className="primary" type="submit" disabled={!editName.trim() || !editCategoryId}><Icon name="check" /><span>Save Changes</span></button></div>
          </form>
        </div>
      )}
      {data.loading || categories.loading ? <div className="spinner" /> : data.error || categories.error ? (
        <div className="error-box card"><p>{data.error || categories.error}</p><button onClick={() => Promise.all([data.load(), categories.load()])}>Try again</button></div>
      ) : groupedCategories.length ? (
        <div className="subcategory-groups">
          {groupedCategories.map(({ category, items }) => (
            <article className="subcategory-group-card card" key={category.id_categories}>
              <div className="subcategory-group-head">
                <span className="subcategory-group-icon"><Icon name="category" /></span>
                <div><h3>{category.name}</h3><p>{items.length} {items.length === 1 ? "subcategory" : "subcategories"}</p></div>
              </div>
              <div className="subcategory-group-body">
                <span className="subcategory-group-label">Subcategories</span>
                {items.length ? (
                  <div className="subcategory-manage-list">
                    {items.map((item) => <div className="subcategory-manage-row" key={item.id_subcategories}><b>{item.namesubcategories}</b><div className="subcategory-row-actions"><button className="detail-action" type="button" onClick={() => setDetailItem(item)}><Icon name="eye" /><span>Detail</span></button><button className="edit-action" type="button" onClick={() => { setEditingItem(item); setEditName(item.namesubcategories || ""); setEditCategoryId(String(item.category?.id_categories ?? item.categoriesId ?? "")); }}><Icon name="edit" /><span>Edit</span></button><button className="danger" type="button" onClick={() => remove(item)}><Icon name="delete" /><span>Delete</span></button></div></div>)}
                  </div>
                ) : <p className="no-subcategories">No subcategories added yet.</p>}
              </div>
            </article>
          ))}
        </div>
      ) : (
        <div className="category-empty card"><span><Icon name="subcategory" /></span><h3>{query ? "No matching subcategories" : "No subcategories yet"}</h3><p>{query ? "Try another search term." : "Add your first subcategory using the button above."}</p></div>
      )}
    </section>
  );
}

function UsersPage() {
  const data = useRows("/admin/users"),
    [q, setQ] = useState("");
  const query = q.trim().toLowerCase();
  const users = data.rows.filter((user) =>
    [user.username, user.email, user.roleName]
      .some((value) => String(value || "").toLowerCase().includes(query)),
  );
  const table = (title: string, list: any[]) => (
    <>
      <h2 className="user-heading">{title}</h2>
      <div className="card table-wrap">
        <table>
          <thead>
            <tr>
              <th>Image</th>
              <th>Name</th>
              <th>Email</th>
              <th>Role</th>
            </tr>
          </thead>
          <tbody>
            {list.map((u) => (
              <tr key={u.id_users}>
                <td>
                  {u.image ? (
                    <img
                      className="avatar-image"
                      src={u.image.startsWith("http") || u.image.startsWith("data:") ? u.image : `data:image/webp;base64,${u.image}`}
                      alt={`${u.username || "User"} profile`}
                    />
                  ) : (
                    <span className="avatar">
                      <Icon name="user" />
                    </span>
                  )}
                </td>
                <td>{u.username || "-"}</td>
                <td>{u.email || "-"}</td>
                <td>{u.roleName || (Number(u.roleId) === 2 ? "Admin" : "User")}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {!list.length && <p className="empty">{query ? "No matching users" : `No ${title.toLowerCase()} accounts found`}</p>}
      </div>
    </>
  );
  return (
    <section>
      <div className="title-row">
        <h2>Manage Users</h2>
        <div className="page-actions">
          <div className="toolbar">
            <input
              className="search"
              placeholder="Search users..."
              aria-label="Search users"
              value={q}
              onChange={(e) => setQ(e.target.value)}
            />
          </div>
        </div>
      </div>
      {data.loading ? (
        <div className="spinner" />
      ) : data.error ? (
        <div className="error-box card"><p>{data.error}</p><button onClick={data.load}>Try again</button></div>
      ) : (
        <>
          {table(
            "Admin",
            users.filter((x) => Number(x.roleId) === 2),
          )}
          {table(
            "User",
            users.filter((x) => Number(x.roleId) === 1),
          )}
        </>
      )}
    </section>
  );
}

function SOSPage() {
  return (
    <ResourcePage
      config={{
        title: "Manage SOS (Save Our Souls)",
        path: "/admin/sos",
        add: true,
        columns: [
          { label: "Name", value: (r) => r.name_sos },
          { label: "Location", value: (r) => r.alamat_sos },
          { label: "Number", value: (r) => r.telepon },
        ],
      }}
    />
  );
}

function FacilitiesPage() {
  return (
    <ResourcePage
      config={{
        title: "Manage Facilities",
        path: "/admin/facility",
        add: true,
        columns: [
          {
            label: "Icon",
            value: (r) =>
              r.icon ? (
                <img
                  className="thumb"
                  src={
                    r.icon.startsWith("http")
                      ? r.icon
                      : `data:image/webp;base64,${r.icon}`
                  }
                  alt=""
                />
              ) : (
                "-"
              ),
          },
          { label: "Name", value: (r) => r.namefacility },
        ],
      }}
    />
  );
}

