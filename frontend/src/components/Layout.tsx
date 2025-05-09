import { useState, useContext } from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import {
  AppBar,
  Box,
  CssBaseline,
  Divider,
  Drawer,
  IconButton,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Toolbar,
  Typography,
  Button,
  useMediaQuery,
  useTheme,
} from '@mui/material';
import {
  Menu as MenuIcon,
  Dashboard as DashboardIcon,
  LocalBar as LocalBarIcon,
  Restaurant as RestaurantIcon,
  ShoppingCart as ShoppingCartIcon,
  People as PeopleIcon,
  Logout as LogoutIcon,
} from '@mui/icons-material';
import { AuthContext } from '../contexts/AuthContext';

const drawerWidth = 240;

interface NavigationItem {
  name: string;
  path: string;
  icon: JSX.Element;
  roles: string[];
}

const navigationItems: NavigationItem[] = [
  { name: 'Dashboard', path: '/', icon: <DashboardIcon />, roles: ['admin', 'bartender'] },
  { name: 'Recipes', path: '/recipes', icon: <LocalBarIcon />, roles: ['admin', 'bartender', 'guest'] },
  { name: 'Ingredients', path: '/ingredients', icon: <RestaurantIcon />, roles: ['admin', 'bartender'] },
  { name: 'Orders', path: '/orders', icon: <ShoppingCartIcon />, roles: ['admin', 'bartender', 'guest'] },
  { name: 'Users', path: '/users', icon: <PeopleIcon />, roles: ['admin'] },
];

const Layout = () => {
  const { user, logout } = useContext(AuthContext);
  const [mobileOpen, setMobileOpen] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  const handleNavigation = (path: string) => {
    navigate(path);
    if (isMobile) {
      setMobileOpen(false);
    }
  };

  // Filter navigation items based on user role
  const filteredNavigationItems = navigationItems.filter((item) => 
    user && item.roles.includes(user.role)
  );

  const drawer = (
    <Box>
      <Toolbar>
        <Typography variant="h6" noWrap component="div">
          Bartender App
        </Typography>
      </Toolbar>
      <Divider />
      <List>
        {filteredNavigationItems.map((item) => (
          <ListItem key={item.path} disablePadding>
            <ListItemButton
              selected={location.pathname === item.path}
              onClick={() => handleNavigation(item.path)}
            >
              <ListItemIcon>{item.icon}</ListItemIcon>
              <ListItemText primary={item.name} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
      <Divider />
      <List>
        <ListItem disablePadding>
          <ListItemButton onClick={logout}>
            <ListItemIcon>
              <LogoutIcon />
            </ListItemIcon>
            <ListItemText primary="Logout" />
          </ListItemButton>
        </ListItem>
      </List>
    </Box>
  );

  return (
    <Box sx={{ display: 'flex' }}>
      <CssBaseline />
      <AppBar
        position="fixed"
        sx={{
          width: { md: `calc(100% - ${drawerWidth}px)` },
          ml: { md: `${drawerWidth}px` },
        }}
      >
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ mr: 2, display: { md: 'none' } }}
          >
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" noWrap component="div" sx={{ flexGrow: 1 }}>
            {/* Show current page title based on path */}
            {filteredNavigationItems.find(item => item.path === location.pathname)?.name || 'Bartender App'}
          </Typography>
          {user && (
            <Box sx={{ display: 'flex', alignItems: 'center' }}>
              <Typography variant="body2" sx={{ mr: 2 }}>
                {user.username} ({user.role})
              </Typography>
              {!isMobile && (
                <Button color="inherit" onClick={logout} startIcon={<LogoutIcon />}>
                  Logout
                </Button>
              )}
            </Box>
          )}
        </Toolbar>
      </AppBar>

      <Box
        component="nav"
        sx={{ width: { md: drawerWidth }, flexShrink: { md: 0 } }}
      >
        {/* Mobile drawer */}
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{
            keepMounted: true, // Better open performance on mobile
          }}
          sx={{
            display: { xs: 'block', md: 'none' },
            '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
          }}
        >
          {drawer}
        </Drawer>

        {/* Desktop drawer */}
        <Drawer
          variant="permanent"
          sx={{
            display: { xs: 'none', md: 'block' },
            '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
          }}
          open
        >
          {drawer}
        </Drawer>
      </Box>

      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: 3,
          width: { md: `calc(100% - ${drawerWidth}px)` },
          minHeight: '100vh',
        }}
      >
        <Toolbar /> {/* Spacer for fixed AppBar */}
        <Outlet /> {/* Render child routes */}
      </Box>
    </Box>
  );
};

export default Layout; 