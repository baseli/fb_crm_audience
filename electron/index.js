const { app, BrowserWindow, ipcMain, dialog, Menu, Tray, nativeImage } = require('electron')
const child_process = require('child_process')

let win = null;
let processInstance = null;
let tray = null;

function createWindow () {
    Menu.setApplicationMenu(null)

    win = new BrowserWindow({
        width: 1280,
        height: 720,
        webPreferences: {
            nodeIntegration: true
        },
    })

    win.loadURL('http://127.0.0.1:19597/web/index.html')
}

app.whenReady().then(() => {
    tray = new Tray(nativeImage.createFromPath(`${__dirname}/icon/logo@2x.png`))
    const menus = Menu.buildFromTemplate([
        {
            label: '打开窗口',
            type: 'normal',
            click: () => {
                if (win.isDestroyed()) {
                    createWindow()
                } else {
                    win.show()
                }
            }
        },
        {
            label: '退出',
            type: 'normal',
            click: () => {
                if (! win.isDestroyed()) {
                    win.close()
                    win.destroy()
                }

                if (processInstance) {
                    if (process.platform === 'win32') {
                        child_process.spawn("taskkill", ["/pid", processInstance.pid, '/f', '/t']);
                    } else {
                        processInstance.kill('SIGKILL')
                    }
                }

                app.quit()
            }
        },
    ])

    createWindow()
    processInstance = child_process.spawn(`${__dirname}/binary/fb_crm_audience`)
    tray.setContextMenu(menus)
})

app.on('quit', () => {
    if (processInstance) {
        if (process.platform === 'win32') {
            child_process.spawn("taskkill", ["/pid", processInstance.pid, '/f', '/t']);
        } else {
            processInstance.kill('SIGKILL')
        }
    }
})

app.on('window-all-closed', () => {
})

app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
        createWindow()
    }
})

ipcMain.on('auth', (ipcEvent, args) => {
    const sub = new BrowserWindow({
      width: 600,
      height: 600,
      parent: win,
    })

    sub.loadURL(`https://www.facebook.com/v8.0/dialog/oauth?client_id=1558910020961842&redirect_uri=https://www.facebook.com/connect/login_success.html&state=${new Date().getTime()}&response_type=token&scope=email,public_profile,ads_management`)
    sub.webContents.on('did-navigate', (event, url) => {
        if (url.startsWith("https://www.facebook.com/connect/login_success.html#")) {
          let params = new URLSearchParams(url.replaceAll('https://www.facebook.com/connect/login_success.html#', ''))
          if (params.has('access_token')) {
              ipcEvent.sender.send('auth', params.get('access_token'))
          }
        }
    })

    ipcEvent.sender.send('auth')
})

ipcMain.on('chooseFile', (event) => {
  dialog.showOpenDialog({
    properties: ['openFile', 'multiSelections'],
    filters: [
      {name: 'Csv', extensions: ['csv']}
    ]
  }).then(res => {
    if (res.canceled === false) {
      event.sender.send('chooseFile', res.filePaths)
    } else {
      event.sender.send('chooseFile')
    }
  }).catch(err => {
    event.sender.send('chooseFile')
  })
})
