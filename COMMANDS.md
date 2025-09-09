# 🏸 SJSU Badminton Bot - Command Reference

The SJSU Badminton Discord Bot supports **both slash commands** (modern Discord feature) and **traditional "!" prefix commands** for maximum compatibility and user convenience.

## 🎯 **Available Commands**

### 🏸 **Court Status Commands**

#### Check Mac Gym Court Availability
- **Slash Command:** `/macgym`
- **Prefix Command:** `!macgym`
- **Description:** Shows current Mac Gym badminton court occupancy
- **Example:** `!macgym`

**Response:** Shows available courts, courts in use, and last updated time.

---

### 📅 **Event Commands**

#### General Badminton Information
- **Slash Command:** `/badminton`
- **Prefix Command:** `!badminton`
- **Description:** Shows general badminton information and available subcommands
- **Example:** `!badminton`

#### List Upcoming Events
- **Slash Command:** `/badminton events [days]`
- **Prefix Command:** `!badminton events [days]`
- **Description:** Lists upcoming badminton events
- **Parameters:**
  - `days` (optional): Number of days to look ahead (default: 7, max: 30)
- **Examples:**
  - `!badminton events` (shows next 7 days)
  - `!badminton events 14` (shows next 14 days)

---

### 🔔 **Alert Commands**

#### Subscribe to Alerts
- **Slash Command:** `/subscribe [threshold]`
- **Prefix Command:** `!subscribe [threshold]`
- **Description:** Subscribe to badminton alerts and notifications
- **Parameters:**
  - `threshold` (optional): Alert when Mac Gym occupancy reaches this level (default: 0)
- **Examples:**
  - `!subscribe` (subscribe to all alerts)
  - `!subscribe 5` (alert when 5+ courts are in use)

#### Unsubscribe from Alerts
- **Slash Command:** `/unsubscribe`
- **Prefix Command:** `!unsubscribe`
- **Description:** Unsubscribe from all badminton alerts
- **Example:** `!unsubscribe`

---

### ℹ️ **Help Commands**

#### Show Help
- **Prefix Command:** `!help`
- **Description:** Shows available commands and usage information
- **Example:** `!help`

**Note:** Slash commands automatically show help through Discord's built-in interface.

---

## 🔄 **Command Types Comparison**

| Feature | Slash Commands (`/`) | Prefix Commands (`!`) |
|---------|---------------------|----------------------|
| **Auto-completion** | ✅ Yes | ❌ No |
| **Parameter validation** | ✅ Yes | ❌ Manual |
| **Modern Discord UI** | ✅ Yes | ❌ No |
| **Backward compatibility** | ❌ Newer Discord | ✅ All versions |
| **Quick typing** | ❌ Slower | ✅ Faster |
| **Mobile friendly** | ✅ Yes | ✅ Yes |

---

## 📱 **Usage Examples**

### **Slash Commands (Modern)**
```
/macgym
/badminton events 14
/subscribe 3
/unsubscribe
```

### **Prefix Commands (Traditional)**
```
!macgym
!badminton events 14
!subscribe 3
!unsubscribe
!help
```

---

## 🎨 **Command Features**

### **Smart Parsing**
- ✅ **Case Insensitive:** `!MACGYM` works the same as `!macgym`
- ✅ **Flexible Spacing:** `!badminton   events   7` handles multiple spaces
- ✅ **Argument Validation:** Invalid numbers default to safe values
- ✅ **Error Handling:** Unknown commands show helpful error messages

### **Rich Responses**
- 🎨 **Embedded Messages:** Beautiful, formatted responses
- 🏷️ **Color Coding:** Different colors for different types of information
- 📊 **Structured Data:** Organized information in fields
- ⏰ **Timestamps:** Shows when data was last updated

### **User-Friendly**
- 👤 **User Mentions:** Alerts mention users by name
- 🔔 **Smart Notifications:** Context-aware alert messages
- 📝 **Clear Instructions:** Helpful error messages and guidance
- 🎯 **Consistent Interface:** Same functionality across both command types

---

## 🚀 **Getting Started**

### **For New Users:**
1. Try slash commands first: `/help` (if available in your Discord version)
2. Use prefix commands: `!help` for a complete command list
3. Start with `!macgym` to check court availability
4. Subscribe to alerts with `!subscribe` for notifications

### **For Power Users:**
- Use `!badminton events 30` to see a full month of events
- Set custom thresholds: `!subscribe 7` to alert when 7+ courts are busy
- Combine commands: Check courts with `!macgym`, then subscribe with `!subscribe 5`

---

## 🔧 **Technical Details**

### **Command Processing**
1. **Message Detection:** Bot detects messages starting with `!`
2. **Parsing:** Splits command and arguments using spaces
3. **Validation:** Checks command exists and arguments are valid
4. **Execution:** Runs appropriate handler function
5. **Response:** Sends formatted response to channel

### **Error Handling**
- **Unknown Commands:** Shows helpful error with suggestion to use `!help`
- **Invalid Arguments:** Gracefully handles bad input with defaults
- **Missing Data:** Shows appropriate "no data available" messages
- **Rate Limiting:** Respects Discord's rate limits

### **Performance**
- **Fast Response:** Commands respond in under 1 second
- **Efficient Parsing:** Minimal CPU usage for command processing
- **Memory Efficient:** No persistent command state stored
- **Scalable:** Handles multiple concurrent commands

---

## 🎾 **Pro Tips**

1. **Use Both Types:** Slash commands for discovery, prefix commands for speed
2. **Set Alerts:** Subscribe with `!subscribe 3` to get notified when courts get busy
3. **Check Events:** Use `!badminton events 14` to plan ahead
4. **Quick Check:** `!macgym` is the fastest way to see current court status
5. **Help Others:** Share `!help` with new users

---

**Happy Badminton! 🏸🤖**

*Both command types work seamlessly together - use whichever you prefer!*
